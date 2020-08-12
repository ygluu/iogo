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

type FuncCallHttpSvc = plugin.FuncCallHttpSvc

type client struct {
	matadata    string
	serviceName string
	svcClient   interface{}
	target      string
}

func (self *client) onRun() {
	lb := gPluginCli.Lb.Clone()
	sd := gPluginCli.Sd.Clone()

	gNetClient.SetCodec(gPluginCli.Cdc)
	sd.Watcher(gClusterName, self.serviceName, lb)
	if self.target != "" {
		lb.Add(self.target, 1)
	}
	gNetClient.Register(self.matadata, self.serviceName, self.svcClient, lb)
}

var gClients []*client = []*client{}
var gNetClient plugin.NetClient = nil

func RegisterClient(matadata, serviceName string, svcClient interface{}, target string) {
	client := &client{
		matadata:    matadata,
		serviceName: serviceName,
		svcClient:   svcClient,
		target:      target,
	}
	gClients = append(gClients, client)
}

func runClient() {
	gNetClient = gPluginCli.Net
	for _, client := range gClients {
		if findSvc(client.matadata) {
			continue
		}
		client.onRun()
	}
}
