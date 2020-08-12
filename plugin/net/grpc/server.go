/*
 *
 * Copyright 2019 The iogo authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

/*
 * contacts: nike: freetoo name: yigui-lu(卢益贵)
 *    wx/qq: 48092788
 *   e-mail: gcode@qq.com
 *      url: https://github.com/freetoo/iogo
 */

package grpc // import "iogo/plugin/net/grpc"

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"sync/atomic"
	"time"

	ggrpc "google.golang.org/grpc"

	"iogo/log"
	"iogo/plugin"
)

type server struct {
	reqCount          int64
	respCount         int64
	ln                GrpcListen
	svr               *GrpcServer
	codec             plugin.Codec
	unaryInterceptor  plugin.FuncUnaryInterceptor
	streamInterceptor plugin.FuncStreamInterceptor
}

func (self *server) Register(matadata, serviceName string, svcServer interface{},
	unaryInterceptor plugin.FuncUnaryInterceptor, streamInterceptor plugin.FuncStreamInterceptor, args ...interface{}) {
	if len(args) == 0 {
		panic("Register: The grpc plug-in needs a description file")
	}
	t := reflect.TypeOf(args[0])
	if t.String() != "*grpc.ServiceDesc" {
		panic("Register: The grpc plug-in needs '*grpc.ServiceDesc'")
	}
	self.unaryInterceptor = unaryInterceptor
	self.streamInterceptor = streamInterceptor
	fileDesc := args[0].(*ggrpc.ServiceDesc)
	self.svr.RegisterService(fileDesc, svcServer)
}

func (self *server) SetCodec(codec plugin.Codec) {
	self.codec = codec
}

func (self *server) GetConcurrentCount() (reqCount int64, respCount int64) {
	req := atomic.LoadInt64(&self.reqCount)
	resp := atomic.LoadInt64(&self.respCount)
	return req, resp
}

func (self *server) Serve() {
	isTrue := true
	checkErr := func() {
		time.Sleep(2 * time.Second)
		if isTrue {
			log.I(LC_GRPCS_READY, "GrpcServer: GrpcServer is ready")
		}
	}
	go checkErr()
	err := self.svr.Serve(self.ln)
	if err != nil {
		isTrue = false
		log.E(LC_GRPCS_LISTENERR, "GrpcServer: Listen is error: %s", err.Error())
		time.Sleep(1 * time.Second)
	}
}

func (self *server) Close() {
	self.svr.Stop()
}

type interceptor struct {
	server *server
}

func (self *interceptor) unaryInterceptor(ctx context.Context, req interface{}, info *ggrpc.UnaryServerInfo, handler ggrpc.UnaryHandler) (resp interface{}, err error) {
	atomic.AddInt64(&self.server.reqCount, 1)
	defer func() {
		if err == nil {
			atomic.AddInt64(&self.server.respCount, 1)
		}
	}()
	if self.server.unaryInterceptor != nil {
		i := &plugin.UnaryServerInfo{FullMethod: info.FullMethod, Server: info.Server}
		resp, err = self.server.unaryInterceptor(ctx, req, i, handler)
	} else {
		resp, err = handler(ctx, req)
	}
	return
}

func (self *interceptor) streamInterceptor(srv interface{}, ss ggrpc.ServerStream, info *ggrpc.StreamServerInfo, handler ggrpc.StreamHandler) (err error) {
	atomic.AddInt64(&self.server.reqCount, 1)
	defer func() {
		if err == nil {
			atomic.AddInt64(&self.server.respCount, 1)
		}
	}()
	if self.server.streamInterceptor != nil {
		i := &plugin.StreamServerInfo{
			FullMethod:     info.FullMethod,
			IsClientStream: info.IsClientStream,
			IsServerStream: info.IsServerStream,
		}
		hdr := func(srv interface{}, stream interface{}) error {
			return handler(srv, stream.(ggrpc.ServerStream))
		}
		err = self.server.streamInterceptor(srv, ss, i, hdr)
	}
	err = handler(srv, ss)
	return
}

func NewServer(newOpt NewOption) plugin.NetServer {
	ic := interceptor{}
	doNewServer := func(opts ...GrpcServerOption) *GrpcServer {
		opts = append(opts, ggrpc.UnaryInterceptor(ic.unaryInterceptor), ggrpc.StreamInterceptor(ic.streamInterceptor))
		ret := ggrpc.NewServer(opts...)
		return ret
	}
	var svr *GrpcServer
	var ln GrpcListen
	if newOpt.Flag() == flag_server {
		md := newOpt.Matadata().([]interface{})
		svr = md[0].(*GrpcServer)
		ln = md[1].(GrpcListen)
	} else if newOpt.Flag() == flag_listen {
		md := newOpt.Matadata().([]interface{})
		ln = md[0].(GrpcListen)
		gsos := md[1].([]GrpcServerOption)
		svr = doNewServer(gsos...)
	} else if newOpt.Flag() == flag_certfilesvr {
		md := newOpt.Matadata().([]interface{})
		fs := md[0].(map[string]string)
		gsos := md[1].([]GrpcServerOption)

		addr := fs["addr"]
		caFile := fs["cafile"]
		certFile := fs["certfile"]
		keyFile := fs["keyfile"]
		creds, err := GetCredsServerOption(certFile, keyFile, caFile)
		if err != nil {
			panic(fmt.Sprintf("NewGrpcServer: create creds: %s", err.Error()))
		}
		if creds != nil {
			gsos = append(gsos, creds)
		}

		svr = doNewServer(gsos...)
		ln, err = net.Listen("tcp", addr)
		if err != nil {
			panic("NewGrpcServer: net.Listen:" + err.Error())
		}
	} else if newOpt.Flag() == flag_addr {
		md := newOpt.Matadata().([]interface{})
		addr := md[0].(string)
		gsos := md[1].([]GrpcServerOption)

		l, err := net.Listen("tcp", addr)
		if err != nil {
			panic("NewGrpcServer: net.Listen error:" + err.Error())
		}
		ln = l
		svr = doNewServer(gsos...)
	} else {
		t := reflect.TypeOf(newOpt.Matadata())
		panic("NewGrpcServer: NewOption error: " + t.String())
	}

	ret := &server{
		reqCount:  0,
		respCount: 0,
		svr:       svr,
		ln:        ln,
	}

	ic.server = ret
	return ret
}
