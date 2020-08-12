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

const strCfgNet = `package iogo

import "recclustername/comm/iogo/grpc"

func RunIogo(servicePort string) {
	iogo.RunIogo(servicePort)
}

func StopIogo() {
	iogo.StopIogo()
}
`

func CreateCfgNet(ns *Names) {
	dir := fmt.Sprintf("%s/%s/comm/iogo", ns.WorkerDir, ns.ClusetNameLower)
	os.MkdirAll(dir, os.ModePerm)
	fn := fmt.Sprintf("%s/iogo.go", dir)
	str := strings.Replace(strCfgNet, clusterNameLower, ns.ClusetNameLower, -1)
	StrToFile(str, fn)
}
