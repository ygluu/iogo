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

package iogoc

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"iogo/iogoc/iogoc/defsrc"
)

func mkNames(workdir, cluster, service string) *defs.Names {
	ret := &defs.Names{
		ClusetNameFUp:    firstUp(cluster),
		ClusetNameLower:  strings.ToLower(cluster),
		ClusetNameUp:     strings.ToUpper(cluster),
		ClusetNameFLow:   firstLower(cluster),
		ServiceNameFUp:   firstUp(service),
		ServiceNameFLow:  firstLower(service),
		ServiceNameLower: strings.ToLower(service),
		ServiceNameUp:    strings.ToUpper(service),
		WorkerDir:        workdir,
		ClusterDir:       fmt.Sprintf("%s/%s", workdir, cluster),
	}
	return ret
}

func domakes(ns *defs.Names) bool {
	defs.CreateProto(ns)
	defs.CreateInterceptor(ns)
	defs.CreateLogCode(ns)
	defs.CreateMethod(ns)
	defs.CreateMain(ns)
	fmt.Println("Service creation completed:", ns.ServiceNameFUp)
	return true
}

func makes(workdir, cluster, service string) bool {
	ns := mkNames(workdir, cluster, service)
	sdir := fmt.Sprintf("%s/%s/svc%s", ns.WorkerDir, ns.ClusetNameLower, ns.ClusetNameLower)
	if PathExists(sdir) {
		fmt.Println("Service dir already exists:", sdir)
		return false
	}
	defs.SetFirstSvc(ns.ClusterDir, false)
	return domakes(ns)
}

func makec(workdir, cluster string) bool {
	ns := mkNames(workdir, cluster, "Demo")
	cdir := fmt.Sprintf("%s/%s", ns.WorkerDir, ns.ClusetNameLower)
	if PathExists(cdir) {
		fmt.Println("Cluster dir already exists:", cdir)
		return false
	}
	defs.SetFirstSvc(ns.ClusterDir, true)
	os.MkdirAll(fmt.Sprintf("%s/proto/src", cdir), os.ModePerm)
	defs.CreateCfgGrpc(ns)
	defs.CreateCfgHttp(ns)
	defs.CreateCfgNet(ns)
	fmt.Println("Cluster creation completed:", ns.ClusetNameFUp)
	domakes(ns)
	defs.CreateDemoMain(ns)
	defs.CreateDemoTest(ns)
	return true
}

func getGoFiles(pathname string) []string {
	ret := []string{}
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		return ret
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := pathname + "/" + fi.Name()
			ret = append(ret, getGoFiles(fullDir)...)
		} else {
			if path.Ext(fi.Name()) != ".go" {
				continue
			}
			if strings.Index(fi.Name(), ".pb.go") != -1 {
				continue
			}
			if strings.Index(fi.Name(), ".ig.go") != -1 {
				continue
			}
			fullName := pathname + "/" + fi.Name()
			ret = append(ret, fullName)
		}
	}
	return ret
}

func execIogoc(protodir string) bool {
	files := getGoFiles(protodir)
	for _, v := range files {
		fmt.Println("Process file:", v)
		writeIogoHttpByGo(v)
	}
	setImportHttp(protodir)
	return false
}

func doWriteIogoHttp(unitName string, methods map[string]*[]string, srcname string, igFileName string) {
	cstr := ""
	nameUp := strings.Replace(path.Base(igFileName), ".ig.go", "", -1)
	nameUp = firstUp(nameUp)
	nameLower := firstLower(nameUp)

	addCli := func(name string) {
		str := `

func (self *%sClient) %s(ctx Context, request *%sRequest) (*%sReply, error) {
	reply := &%sReply{}
	return reply, self.CallMethod("/proto.%s/%s", ctx, request, reply)
}`
		cstr += fmt.Sprintf(str, nameLower, name, name, name, name, nameUp, name)
	}

	sstr := ""
	hstr := ""
	istr := ""
	addSvr := func(name string) {
		str := `
	%s(ctx Context, request *%sRequest) (*%sReply, error)`
		sstr += fmt.Sprintf(str, name, name, name)

		str = `

func %sHandler%s(svr interface{}, ctx Context, request interface{}) (interface{}, error) {
	return svr.(%sServer).%s(ctx, request.(*%sRequest))
}`
		hstr += fmt.Sprintf(str, strings.ToLower(name), nameUp, nameUp, name, name)

		str = `
	"%s": http.MethodInfo{
		HandlerMethod:  %sHandler%s,
		FuncNewRequest: func() interface{} { return &%sRequest{} },
	},`
		istr += fmt.Sprintf(str, name, strings.ToLower(name), nameUp, name)
	}

	for k, v := range methods {
		req := (*v)[0]
		res := (*v)[1]
		if (req != "") && (res == "") {
			fmt.Printf("Error: The message %sRequest not define Reply in file %s\n", k, srcname)
			continue
		}
		if (req == "") && (res != "") {
			fmt.Printf("Error: The message %sReply not define Request in file %s\n", k, srcname)
			continue
		}
		addCli(k)
		addSvr(k)
	}

	strfile := `package proto

import "iogo"
import "iogo/unit/md"
import "iogo/plugin/net/http"

import (
	"context"
)

type Context = context.Context

func TODO() context.Context {
	return context.TODO()
}

func Background() context.Context {
	return context.Background()
}

func TokenContext(token string) context.Context {
	ctx := TODO()
	ctx = md.AppendToOutgoingContext(ctx, "set-token", token)
	return ctx
}

type %sClient struct {
	CallMethod iogo.FuncCallHttpSvc
}%s

var Call%s *%sClient = &%sClient{}

func onStart() {

}

func init() {
	iogo.RegisterFuncOnStart(onStart)
	
	iogo.RegisterClient("%s", "proto.%s", Call%s, "")
}

type %sServer interface {%s
}%s

var methodInfos%s http.MethodInfos = http.MethodInfos{%s
}

func Register%sServer(srv %sServer) {
	iogo.RegisterServer("%s", "proto.%s", srv, methodInfos%s)
}
`
	strfile = fmt.Sprintf(strfile, nameLower, cstr, nameUp, nameLower, nameLower,
		unitName, nameUp, nameUp, nameUp, sstr, hstr, nameUp, istr, nameUp, nameUp, unitName, nameUp, nameUp)
	StrToFile(strfile, igFileName)
	return
}

func setImportHttp(protodir string) {
	fn := strings.Replace(protodir, "proto", "comm/iogo/iogo.go", -1)
	lines := fileToLine(fn)
	if len(lines) == 0 {
		return
	}
	str := ""
	for _, line := range lines {
		str += line
	}
	str = strings.Replace(str, "/comm/iogo/grpc", "/comm/iogo/http", -1)
	StrToFile(str, fn)
}

func writeIogoHttpByGo(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	methods := make(map[string]*[]string)
	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		str := strings.Trim(line, "\n")
		str = strings.Trim(str, "\r")

		msg := getBetweenStr(str, "type ", "Request struct {")
		if msg != "" {
			method := methods[msg]
			if method == nil {
				m := make([]string, 2)
				method = &m
				methods[msg] = &m
			}
			(*method)[0] = msg
			continue
		}
		msg = getBetweenStr(str, "type ", "Reply struct {")
		if msg != "" {
			method := methods[msg]
			if method == nil {
				m := make([]string, 2)
				method = &m
				methods[msg] = &m
			}
			(*method)[1] = msg
		}
	}

	pbgoDir := path.Dir(filename)
	dirs := strings.Split(pbgoDir, "/")
	unitName := strings.Replace(path.Base(filename), ".proto", "", -1)
	unitName = fmt.Sprintf("%s.%s.proto", dirs[len(dirs)-1], unitName)
	//unitName = strings.ToLower(unitName)

	doWriteIogoHttp(unitName, methods, filename, toIgFileNmae(filename))
}

func toIgFileNmae(filename string) string {
	dir := path.Dir(filename)
	name := path.Base(filename)
	arr := strings.Split(name, ".")
	return fmt.Sprintf("%s/%s.ig.go", dir, arr[0])
}
