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

package defs

import (
	"fmt"
	"os"
	"strings"
)

const strMethod1 = `package server

import (
	"recclustername/proto/recservicename"
)

type recServiceNameServer struct {
}

func init() {
	proto.RegisterRecServiceNameServer(&recServiceNameServer{})
}

func (self *recServiceNameServer) Hellow(ctx proto.Context, request *proto.HellowRequest) (*proto.HellowReply, error) {
	reply := &proto.HellowReply{
		Answer: "Welcome to here",
	}
	return reply, nil
}
`

const strMethod2 = `package server

import (
	"recclustername/proto/recservicename"
)

type recServiceNameServer struct {
}

func init() {
	proto.RegisterRecServiceNameServer(&recServiceNameServer{})
}
`

func CreateMethod(ns *Names) {
	dir := fmt.Sprintf("%s/%s/svc%s/server", ns.WorkerDir, ns.ClusetNameLower, ns.ServiceNameLower)
	os.MkdirAll(dir, os.ModePerm)
	fn := fmt.Sprintf("%s/method.go", dir)
	str := strMethod2
	if isFirstSvr {
		str = strMethod1
	}
	str = strings.Replace(str, clusterNameLower, ns.ClusetNameLower, -1)
	str = strings.Replace(str, serviceNameLower, ns.ServiceNameLower, -1)
	str = strings.Replace(str, clusterNameFUp, ns.ClusetNameFUp, -1)
	str = strings.Replace(str, serviceNameFUp, ns.ServiceNameFUp, -1)
	str = strings.Replace(str, serviceNameFLow, ns.ServiceNameFLow, -1)
	StrToFile(str, fn)
}
