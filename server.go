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

package iogo // import "iogo"

import (
	"iogo/plugin"
)

type server struct {
	serviceName string
	matadata    string
	svr         interface{}
	args        []interface{}
}

var gFuncGetWeight plugin.FuncGetWeight = func() int { return 1 }

func SetFuncGetWeight(funcGetWeight plugin.FuncGetWeight) {
	gFuncGetWeight = funcGetWeight
}

func (self *server) onRun() {
	sd := gPluginSvr.Sd.Clone()
	sd.Register(gClusterName, self.serviceName, gServiceAddr, gFuncGetWeight)
	gNetServer.Register(self.matadata, self.serviceName, self.svr, gUnaryInterceptor, gStreamInterceptor, self.args...)
}

var gRegServers []*server = []*server{}

func RegisterServer(matadata, serviceName string, svr interface{}, args ...interface{}) {
	server := &server{
		serviceName: serviceName,
		matadata:    matadata,
		svr:         svr,
		args:        args,
	}
	gRegServers = append(gRegServers, server)
}

func findSvc(matadata string) bool {
	for _, server := range gRegServers {
		if matadata == server.matadata {
			return true
		}
	}
	return false
}

func runServer() {
	if !hasServer() && len(gRegServers) > 0 {
		panic("IogoServer: The server plug-in is not set")
	}
	if hasServer() && len(gRegServers) == 0 {
		panic("IogoServer: Please use RegisterServer to register the server")
	}
	if hasServer() {
		gNetServer.SetCodec(gPluginSvr.Cdc)
		for _, server := range gRegServers {
			server.onRun()
		}
	}
}

type UnaryServerInfo = plugin.UnaryServerInfo
type UnaryHandler = plugin.UnaryHandler
type FuncUnaryInterceptor = plugin.FuncUnaryInterceptor

type StreamServerInfo = plugin.StreamServerInfo
type StreamHandler = plugin.StreamHandler
type FuncStreamInterceptor = plugin.FuncStreamInterceptor

var gUnaryInterceptor FuncUnaryInterceptor = nil
var gStreamInterceptor FuncStreamInterceptor = nil

func RegisterIntercaptors(unaryInterceptor FuncUnaryInterceptor, streamInterceptor FuncStreamInterceptor) {
	if gUnaryInterceptor != nil {
		panic("Interceptor can only be registered once")
	}
	gUnaryInterceptor = unaryInterceptor
	gStreamInterceptor = streamInterceptor
}
