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
	"context"

	"iogo/log"
	"iogo/plugin"
)

var (
	Log = log.Get()
)

type Context = context.Context

type pluginClient struct {
	Lb  plugin.LoadBalan
	Net plugin.NetClient
	Sd  plugin.DiscClient
	Cdc plugin.Codec
}

func NewPluginClient() *pluginClient {
	return &pluginClient{}
}

type pluginServer struct {
	Net plugin.NetServer
	Sd  plugin.DiscServer
	Cdc plugin.Codec
}

func NewPluginServer() *pluginServer {
	return &pluginServer{}
}

var gPluginSvr *pluginServer = nil
var gPluginCli *pluginClient = nil

func hasServer() bool {
	return gPluginSvr != nil && gPluginSvr.Net != nil &&
		gPluginSvr.Sd != nil && gPluginSvr.Cdc != nil
}

func hasClient() bool {
	return gPluginCli != nil && gPluginCli.Net != nil &&
		gPluginCli.Sd != nil && gPluginCli.Cdc != nil && gPluginCli.Lb != nil
}

func SetPlugin(pluginSvr *pluginServer, pluginCli *pluginClient) {
	gPluginSvr = pluginSvr
	gPluginCli = pluginCli
}

var funOnStarts []func() = []func(){}

func RegisterFuncOnStart(onStart func()) {
	funOnStarts = append(funOnStarts, onStart)
}

var gClusterName string = ""
var gServiceAddr string = ""
var gNetServer plugin.NetServer = nil

func Run(clusterName string, serviceAddr string) {
	gClusterName = clusterName
	gServiceAddr = serviceAddr

	if hasServer() {
		if gServiceAddr == "" {
			panic("IogoServer: Service address not set, Check the settings: service port and subnet mask, or service address")
		}
		gNetServer = gPluginSvr.Net
		gNetServer.SetCodec(gPluginSvr.Cdc)
	}

	log.Get().SetKey(Runtime.GetServiceKey())
	runClient()
	runServer()

	gPluginCli = nil
	gPluginSvr = nil
	startRuntime()
	Log = log.Get()

	for _, funcOnStart := range funOnStarts {
		funcOnStart()
	}

	if gNetServer != nil {
		gNetServer.Serve()
	}
}

func Stop() {
	gNetServer.Close()
}
