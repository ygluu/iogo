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

const strMain = `package main

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"recclustername/comm/iogo"
	_ "recclustername/svcrecservicename/server"
)

const defServicePort = "28666"

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	isError := true
	waitStop := func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
		go func() {
			<-ch
			isError = false
			iogo.StopIogo()
			os.Exit(1)
		}()
	}
	waitStop()

	iogo.RunIogo(defServicePort)
	if isError {
		iogo.StopIogo()
	}
}
`

func CreateMain(ns *Names) {
	dir := fmt.Sprintf("%s/%s/svc%s", ns.WorkerDir, ns.ClusetNameLower, ns.ServiceNameLower)
	os.MkdirAll(dir, os.ModePerm)
	fn := fmt.Sprintf("%s/main.go", dir)
	str := strings.Replace(strMain, clusterNameLower, ns.ClusetNameLower, -1)
	str = strings.Replace(str, serviceNameLower, ns.ServiceNameLower, -1)
	str = strings.Replace(str, "28666", fmt.Sprintf("%d", 28666+cIndex), -1)
	StrToFile(str, fn)
}
