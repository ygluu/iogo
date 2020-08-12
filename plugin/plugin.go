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

package plugin

import "context"

const LC_LOG_MODULE = 200
const LC_ETCD_MODULE = 300
const LC_ZK_MODULE = 400
const LC_HTTP_MODULE = 500
const LC_GRPC_MODULE = 600
const LC_LB_MODULE = 700

type FuncGetWeight func() int

type FuncCallHttpSvc func(fullMethod string, ctx context.Context, request interface{}, response interface{}) error

type DiscClient interface {
	Clone() DiscClient
	Watcher(clusterName string, serviceName string, loadBalan LoadBalan)
}

type DiscServer interface {
	Clone() DiscServer
	Register(clusterName string, serviceName string, serviceAddr string, funcGetWeight FuncGetWeight)
}

type Codec interface {
	Name() string
	String() string
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

type NetClient interface {
	SetCodec(codec Codec)
	Register(matadata, serviceName string, svcClient interface{}, loadBalan LoadBalan)
}

type NetServer interface {
	GetConcurrentCount() (reqCount int64, respCount int64)
	SetCodec(codec Codec)
	Register(matadata, serviceName string, svcServer interface{},
		unaryInterceptor FuncUnaryInterceptor, streamInterceptor FuncStreamInterceptor, args ...interface{})
	Serve()
	Close()
}

type Address struct {
	Addr     string
	Metadata interface{}
}

type LoadBalan interface {
	Clone() LoadBalan
	Get(key string) string
	Add(addr string, weight int)
	Del(addr string)
	Up(addr string, metadata interface{})
	Down(addr string)
	Notify() <-chan int
	Address() []Address
}

type Logger interface {
	I(code int, format string, v ...interface{})
	W(code int, format string, v ...interface{})
	E(code int, format string, v ...interface{})
	D(code int, format string, v ...interface{})
	T(code int, format string, v ...interface{})
	SetKey(key string)
}

type UnaryServerInfo struct {
	FullMethod string
	Server     interface{}
}
type UnaryHandler = func(ctx context.Context, req interface{}) (reply interface{}, err error)
type FuncUnaryInterceptor = func(ctx context.Context, req interface{}, info *UnaryServerInfo, handler UnaryHandler) (resp interface{}, err error)

type StreamServerInfo struct {
	FullMethod     string
	IsClientStream bool
	IsServerStream bool
}
type StreamHandler = func(srv interface{}, stream interface{}) error
type FuncStreamInterceptor = func(srv interface{}, ss interface{}, info *StreamServerInfo, handler StreamHandler) error
