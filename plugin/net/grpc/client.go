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
	"reflect"
	"time"
	"unsafe"

	ggrpc "google.golang.org/grpc"

	"iogo/log"
	"iogo/plugin"
	"iogo/unit"
)

type myToken struct {
	md map[string]string
}

func (self *myToken) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return self.md, nil
}

func (self *myToken) RequireTransportSecurity() bool {
	return true
}

type svcClientInfo struct {
	owner       *client
	serviceName string
	matadata    string
	target      string
	lb          ggrpc.Balancer
	codec       plugin.Codec
	conn        *ggrpc.ClientConn
}

type client struct {
	gdos    []ggrpc.DialOption
	codec   plugin.Codec
	svccs   []*svcClientInfo
	isFirst bool
}

func (self *client) SetCodec(codec plugin.Codec) {
	self.codec = codec
}

func (self *client) newConnect(lb ggrpc.Balancer, opts ...ggrpc.DialOption) *ggrpc.ClientConn {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	opts = append(opts, ggrpc.WithBalancer(lb))
	opts = append(opts, self.gdos...)
	conn, err := ggrpc.DialContext(ctx, "", opts...)
	if err != nil {
		panic(fmt.Sprintf("grpc.DialContext: %s", err.Error()))
	}
	return conn
}

func (self *client) Register(matadata, serviceName string, svcClient interface{}, loadBalan plugin.LoadBalan) {
	if !unit.IsStructPoint(svcClient) {
		panic("RegisterClient: Parameter svcClient must be a structure pointer")
	}
	var svcPConn **ggrpc.ClientConn = nil
	v := reflect.ValueOf(svcClient).Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := t.Field(i)
		if f.Type.String() == "*grpc.ClientConn" {
			svcPConn = (**ggrpc.ClientConn)(unsafe.Pointer(v.Field(i).Addr().Pointer()))
		}
	}
	if svcPConn == nil {
		panic(fmt.Sprintf("Grpc plug-ins cannot use other source code: %s", t.String()))
	}

	lb := NewBalaner(loadBalan)
	svcConn := self.newConnect(lb)
	*svcPConn = svcConn

	svcc := &svcClientInfo{
		owner:       self,
		serviceName: serviceName,
		matadata:    matadata,
		lb:          lb,
		codec:       self.codec,
		conn:        svcConn,
	}
	self.svccs = append(self.svccs, svcc)
	if self.isFirst {
		self.isFirst = false
		log.I(LC_GRPCC_READY, "GrpcClient: Client is ready")
	}
}

func NewClient(newOpt NewOption) plugin.NetClient {
	var gdos []GrpcDailOption
	if newOpt.Flag() == flag_certfilecli {
		md := newOpt.Matadata().([]interface{})
		fs := md[0].(map[string]string)
		gdos = md[1].([]GrpcDailOption)

		certsvrname := fs["certsvrname"]
		caFile := fs["cafile"]
		certFile := fs["certfile"]
		keyFile := fs["keyfile"]
		creds, err := GetCredsClientOption(certFile, keyFile, caFile, certsvrname)
		if err != nil {
			panic(fmt.Sprintf("NewGrpcClient: create creds: %s", err.Error()))
		}
		if creds == nil {
			gdos = append(gdos, ggrpc.WithInsecure())
		}
		if creds != nil {
			gdos = append(gdos, creds)
		}
	} else if newOpt.Flag() == flag_default {
		gdos = newOpt.Matadata().([]ggrpc.DialOption)
		gdos = append(gdos, ggrpc.WithInsecure())
	} else {
		t := reflect.TypeOf(newOpt.Matadata())
		panic("NewGrpcServer: NewOption error: " + t.String())
	}
	return &client{gdos: gdos, isFirst: true}
}
