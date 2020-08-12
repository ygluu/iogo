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
	"io/ioutil"
	"os"
	"strings"
)

func StrToFile(s string, name string) {
	file, err := os.Create(name)
	if err != nil {
		fmt.Println("File is exists: ", name)
		return
	}
	defer file.Close()
	file.WriteString(s)
}

const (
	serviceNameFUp   = "RecServiceName"
	serviceNameFLow  = "recServiceName"
	serviceNameUp    = "RECSERVICENAME"
	serviceNameLower = "recservicename"
	clusterNameFUp   = "RecClusterName"
	clusterNameFLow  = "recClusterName"
	clusterNameUp    = "RECCLUSTERNAME"
	clusterNameLower = "recclustername"
)

type Names struct {
	WorkerDir        string
	ServiceNameFUp   string
	ServiceNameFLow  string
	ServiceNameUp    string
	ServiceNameLower string
	ClusetNameFUp    string
	ClusetNameFLow   string
	ClusetNameUp     string
	ClusetNameLower  string
	ClusterDir       string
}

var isFirstSvr bool = true
var cIndex int = 0

func SetFirstSvc(clusterDir string, v bool) {
	isFirstSvr = v
	cIndex = getClusterCount(clusterDir)
}

func getClusterCount(p string) int {
	dir, err := ioutil.ReadDir(p)
	if err != nil {
		return 0
	}
	ret := 0
	for _, fi := range dir {
		if fi.IsDir() && (strings.Index(fi.Name(), "svc") == 0) {
			ret++
		}
	}
	return ret
}
