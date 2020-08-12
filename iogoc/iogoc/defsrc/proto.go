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

const strProto1 = `syntax = "proto3";

option objc_class_prefix = "HLW";

package proto;

message HellowRequest {
    string Name = 1;
}

message HellowReply {
    string Answer = 2;
}

service RecServiceName {
	rpc Hellow(HellowRequest) returns (HellowReply) {	
	}
}
`

const strProto2 = `syntax = "proto3";

option objc_class_prefix = "HLW";

package proto;

service RecServiceName {
}
`

func CreateProto(ns *Names) {
	dir := fmt.Sprintf("%s/%s/proto/src/%s", ns.WorkerDir, ns.ClusetNameLower, ns.ServiceNameLower)
	os.MkdirAll(dir, os.ModePerm)
	fn := fmt.Sprintf("%s/%s.proto", dir, ns.ServiceNameFUp)
	str := strProto2
	if isFirstSvr {
		str = strProto1
	}
	str = strings.Replace(str, clusterNameFUp, ns.ClusetNameFUp, -1)
	str = strings.Replace(str, serviceNameFUp, ns.ServiceNameFUp, -1)
	StrToFile(str, fn)
}
