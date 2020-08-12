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

const strInterceptor = `package server

import (
	"fmt"

	"iogo"
	"iogo/unit/md"
)

func init() {
	iogo.RegisterIntercaptors(unaryIntercaptor, streamInterceptor)
}

func streamInterceptor(srv interface{}, ss interface{}, info *iogo.StreamServerInfo, handler iogo.StreamHandler) error {
	return handler(srv, ss)
}

func unaryIntercaptor(ctx iogo.Context, req interface{}, info *iogo.UnaryServerInfo,
	handler iogo.UnaryHandler) (resp interface{}, err error) {
	isTrue := true
	flag := fmt.Sprintf("%p", info)
	if isTrue {
		inh, _ := md.FromIncomingContext(ctx)
		iogo.Log.I(LC_RECSERVICENAME_CALLBEFOR, "UnaryIntercaptorBefor:%s call: %s request: %+v, in-header: %+v", flag, info.FullMethod, req, inh)
	}
	resp, err = handler(ctx, req)
	if isTrue {
		if err != nil {
			iogo.Log.I(LC_RECSERVICENAME_CALLAFTER, "UnaryIntercaptorAfter:%s call: %s error: %v", flag, info.FullMethod, err)
		} else {
			iogo.Log.I(LC_RECSERVICENAME_TOKENINFO, "UnaryIntercaptorAfter:%s call: %s reply: %+v", flag, info.FullMethod, resp)
		}
	}
	return
}
`

func CreateInterceptor(ns *Names) {
	dir := fmt.Sprintf("%s/%s/svc%s/server", ns.WorkerDir, ns.ClusetNameLower, ns.ServiceNameLower)
	os.MkdirAll(dir, os.ModePerm)
	fn := fmt.Sprintf("%s/interceptor.go", dir)
	str := strings.Replace(strInterceptor, clusterNameUp, ns.ClusetNameUp, -1)
	StrToFile(str, fn)
}
