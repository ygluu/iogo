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

const strDemoTest = `package main

import (
	"fmt"
	"runtime"
	"time"

	"recclustername/proto/recservicename"
)

func init() {

}

func testStart() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	time.Sleep(1 * time.Second)

	request := &proto.HellowRequest{
		Name: "My name is iogo",
	}
	fmt.Printf("Request: %+v\n", request)
	reply, err := proto.CallDemo.Hellow(proto.TODO(), request)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Reply: %+v\n", reply)
	}
}
`

func CreateDemoTest(ns *Names) {
	dir := fmt.Sprintf("%s/%s/%sclient", ns.WorkerDir, ns.ClusetNameLower, ns.ServiceNameLower)
	os.MkdirAll(dir, os.ModePerm)
	fn := fmt.Sprintf("%s/test.go", dir)
	str := strings.Replace(strDemoTest, clusterNameLower, ns.ClusetNameLower, -1)
	str = strings.Replace(str, serviceNameLower, ns.ServiceNameLower, -1)
	StrToFile(str, fn)
}
