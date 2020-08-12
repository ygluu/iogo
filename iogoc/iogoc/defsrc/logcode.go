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

const strLogCode = `package server

const LC_RECSERVICENAME_CALLBEFOR = REC2000
const LC_RECSERVICENAME_CALLAFTER = REC2001
const LC_RECSERVICENAME_TOKENINFO = REC2002
`

func CreateLogCode(ns *Names) {
	dir := fmt.Sprintf("%s/%s/svc%s/server", ns.WorkerDir, ns.ClusetNameLower, ns.ServiceNameLower)
	os.MkdirAll(dir, os.ModePerm)
	fn := fmt.Sprintf("%s/logcode.go", dir)
	str := strings.Replace(strLogCode, clusterNameUp, ns.ClusetNameUp, -1)
	str = strings.Replace(str, "REC2", fmt.Sprintf("%d", (cIndex+1)*1000), -1)
	StrToFile(str, fn)
}
