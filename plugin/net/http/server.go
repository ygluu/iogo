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

package http // import "iogo/plugin/net/http"

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"strings"
	"sync/atomic"
	"time"

	"iogo/log"
	"iogo/plugin"
	"iogo/plugin/net/http/router"
	"iogo/unit/md"
	mytls "iogo/unit/tls"
)

type FuncNewRequest = func() interface{}
type HandlerMethod = func(svr interface{}, ctx context.Context, request interface{}) (interface{}, error)
type MethodInfo struct {
	FuncNewRequest FuncNewRequest
	HandlerMethod  HandlerMethod
}
type MethodInfos map[string]MethodInfo

type svcMethod struct {
	owner          interface{}
	server         *server
	funcNewRequest FuncNewRequest
	handlerMethod  HandlerMethod
}

type server struct {
	svr         *HttpServer
	ln          HttpListen
	reqCount    int64
	respCount   int64
	isHttp      bool
	router      *router.Router
	codec       plugin.Codec
	interceptor plugin.FuncUnaryInterceptor
}

func (self *server) doRegService(matadata, serviceName string, svcSvr interface{}, minfos MethodInfos) {
	for k, v := range minfos {
		method := &svcMethod{
			server:         self,
			owner:          svcSvr,
			funcNewRequest: v.FuncNewRequest,
			handlerMethod:  v.HandlerMethod,
		}
		self.router.AddMethod(serviceName, k, method)
	}
}

func (self *svcMethod) docall(ctx context.Context, req interface{}) (reply interface{}, err error) {
	return self.handlerMethod(self.owner, ctx, req)
}

func (self *server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	atomic.AddInt64(&self.reqCount, 1)

	reqBody, err := ioutil.ReadAll(req.Body)
	setError := func(format string, args ...interface{}) {
		estr := fmt.Sprintf(format, args...)
		w.Header().Set(IOGO_ERROR_NAME, IOGO_ERROR_VALUE)
		w.Write([]byte(estr))
	}
	if err != nil {
		setError("URL: %s Read body error: %s", err.Error())
		return
	}

	fullMethod := strings.Trim(req.URL.Path, "/")
	urls := strings.Split(fullMethod, "/")
	l := len(urls)
	if l < 2 {
		setError("URL errror: %s", req.URL.Path)
		return
	}
	methodName := urls[len(urls)-1]
	serviceName := urls[len(urls)-2]

	retm := self.router.GetMethod(serviceName, methodName)
	if retm == nil {
		setError("Not method '%s'", fullMethod)
		return
	}
	method := retm.(*svcMethod)

	request := method.funcNewRequest()
	codec := self.codec
	err = codec.Unmarshal(reqBody, request)
	if err != nil {
		setError("Unmarshal of codec error: %s", err.Error())
		return
	}
	inh := md.MD{}
	for k, v := range req.Header {
		k := strings.ToLower(k)
		inh[k] = v
	}
	ctx := md.NewIncomingContext(req.Context(), inh)
	var response interface{}
	if self.interceptor != nil {
		info := &plugin.UnaryServerInfo{
			FullMethod: req.URL.Path,
			Server:     method.owner,
		}
		response, err = self.interceptor(ctx, request, info, method.docall)
		if err != nil {
			setError(err.Error())
			return
		}
	} else {
		response, err = method.docall(ctx, request)
		if err != nil {
			setError(err.Error())
			return
		}
	}
	if err != nil {
		setError("Unmarshal of codec error: %s", err.Error())
		return
	}
	respBody, err := codec.Marshal(response)
	if err == nil {
		w.Write(respBody)
		atomic.AddInt64(&self.respCount, 1)
	}
}

func (self *server) GetConcurrentCount() (reqCount int64, respCount int64) {
	req := atomic.LoadInt64(&self.reqCount)
	resp := atomic.LoadInt64(&self.respCount)
	return req, resp
}

func (self *server) SetCodec(codec plugin.Codec) {
	self.codec = codec
}

func (self *server) Register(matadata, serviceName string, svcSvr interface{},
	unaryInterceptor plugin.FuncUnaryInterceptor, streamInterceptor plugin.FuncStreamInterceptor, args ...interface{}) {
	checkMInfos := func() bool {
		for _, arg := range args {
			if arg == nil {
				continue
			}
			t := reflect.TypeOf(arg)
			if t.String() == "http.MethodInfos" {
				return true
			}
		}
		return false
	}
	if (len(args) == 0) || (!checkMInfos()) {
		panic("Register: The http plug-in needs a MethodInfos")
	}
	self.interceptor = unaryInterceptor
	self.doRegService(matadata, serviceName, svcSvr, args[0].(MethodInfos))
}

func (self *server) Serve() {
	isTrue := true
	checkErr := func() {
		time.Sleep(2 * time.Second)
		if isTrue {
			log.I(LC_HTTPS_READY, "HttpServer: Http server is ready")
		}
	}

	self.svr.Handler = self
	go checkErr()
	err := self.svr.Serve(self.ln)
	if err != nil {
		isTrue = false
		log.E(LC_HTTPS_LISTENERR, "HttpServer: Listen '%s' is error: %s", self.svr.Addr, err.Error())
		time.Sleep(1 * time.Second)
	}
}

func (self *server) Close() {
	self.svr.Close()
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func NewServer(newOpt NewOption) plugin.NetServer {
	var svr *HttpServer
	var ln HttpListen
	if newOpt.Flag() == flag_server {
		svr := newOpt.Matadata().(*HttpServer)
		if svr.Addr == "" {
			panic("NewHttpServer: Server addr error")
		}
		listen, err := net.Listen("tcp", svr.Addr)
		if err != nil {
			panic("NewHttpServer: net.Listen error:" + err.Error())
		}
		if svr.TLSConfig != nil {
			ln = tls.NewListener(listen, svr.TLSConfig)
		} else {
			ln = listen
		}
	} else if newOpt.Flag() == flag_listen {
		svr = &HttpServer{}
		ln = newOpt.Matadata().(HttpListen)
	} else if newOpt.Flag() == flag_certfilesvr {
		md := newOpt.Matadata().(map[string]string)
		addr := md["addr"]
		caFile := md["cafile"]
		certFile := md["certfile"]
		keyFile := md["keyfile"]
		var err error
		ln, err = net.Listen("tcp", addr)
		if err != nil {
			panic("NewHttpServer: net.Listen error:" + err.Error())
		}
		if certFile != "" {
			cfg, err := mytls.NewTlsServer(certFile, keyFile, caFile)
			if err != nil {
				panic("NewHttpServer: cert file error:" + err.Error())
			}
			l := tcpKeepAliveListener{ln.(*net.TCPListener)}
			ln = tls.NewListener(l, cfg)
		}
		svr = &HttpServer{}
	} else if newOpt.Flag() == flag_addr {
		addr := newOpt.Matadata().(string)
		var err error
		ln, err = net.Listen("tcp", addr)
		if err != nil {
			panic("NewHttpServer: net.Listen error:" + err.Error())
		}
		svr = &HttpServer{}
	} else {
		t := reflect.TypeOf(newOpt.Matadata())
		panic("NewHttpServer: NewOption error: " + t.String())
	}

	return &server{
		svr:       svr,
		ln:        ln,
		reqCount:  0,
		respCount: 0,
		router:    router.NewRouter(),
	}
}
