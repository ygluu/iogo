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
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

func getMsgFullName(s string) string {
	strs := strings.Split(s, " ")
	ret := []string{}
	for _, str := range strs {
		if str != "" {
			ret = append(ret, str)
		}
	}
	if (len(ret) == 3) && (ret[0] == "message") && (ret[2] == "{") {
		return ret[1]
	}
	return ""
}

func firstUp(s string) string {
	if s == "" {
		return ""
	}
	datas := []byte(s)
	datas[0] = ([]byte(strings.ToUpper(string(datas[0]))))[0]
	return string(datas)
}

func firstLower(s string) string {
	if s == "" {
		return ""
	}
	datas := []byte(s)
	datas[0] = ([]byte(strings.ToLower(string(datas[0]))))[0]
	return string(datas)
}

func getMsgName(s string, ext string) string {
	index := strings.Index(s, ext)
	if index != len(s)-len(ext) {
		return ""
	}
	ret := s[0:index]
	return ret
}

func writeIogoHttp(unitName, srcname, pbgoName string) {
	f, err := os.Open(srcname)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	count := 0
	methods := make(map[string]*[]string)
	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		count++
		if strings.Index(line, "message") != 0 {
			continue
		}
		str := strings.Trim(line, "\n")
		str = strings.Trim(str, "\r")
		msgname := getMsgFullName(str)
		if msg := getMsgName(msgname, "Request"); msg != "" {
			method := methods[msg]
			if method == nil {
				m := make([]string, 2)
				method = &m
				methods[msg] = &m
			}
			(*method)[0] = msg
		}
		if msg := getMsgName(msgname, "Reply"); msg != "" {
			method := methods[msg]
			if method == nil {
				m := make([]string, 2)
				method = &m
				methods[msg] = &m
			}
			(*method)[1] = msg
		}
	}

	doWriteIogoHttp(unitName, methods, srcname, toIgFileNmae(pbgoName))
}

func getProtoFiles(pathname string) []string {
	ret := []string{}
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		return ret
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := pathname + "/" + fi.Name()
			ret = append(ret, getProtoFiles(fullDir)...)
		} else {
			if path.Ext(fi.Name()) != ".proto" {
				continue
			}
			fullName := pathname + "/" + fi.Name()
			ret = append(ret, fullName)
		}
	}
	return ret
}

func execCmd(scmd string) bool {
	fmt.Println("exec protoc command:")
	fmt.Println(scmd)
	args := strings.Split(scmd, " ")
	cmd := exec.Command(args[0], args[1:]...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("\nprotoc error:", fmt.Sprint(err))
		fmt.Println(stderr.String())
		return false
	}
	if len(out.String()) > 0 {
		fmt.Println("\nproto out:\n" + out.String())
	}
	return true
}

func fileToLine(filename string) (ret []string) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		ret = append(ret, line)
	}
	return ret
}

func getBetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		return ""
	}
	str = string([]byte(str)[n+len(start):])
	m := strings.Index(str, end)
	if m == -1 {
		return ""
	}
	str = string([]byte(str)[:m])
	return str
}

func checkNewCli(unitName string, lines []string, i int, pgstr *string, varstring *string, initstr *string) int {
	if len(lines)-i < 3 {
		return i
	}

	s1 := lines[i]
	s2 := lines[i+1]
	s3 := lines[i+2]

	cliName := getBetweenStr(s1, "func New", "Client(cc *grpc.ClientConn) ")
	cliType := getBetweenStr(s2, "	return &", "Client{cc}")
	cliIface := getBetweenStr(s1, "Client(cc *grpc.ClientConn) ", " {")
	if cliName == "" || cliType == "" || cliIface == "" {
		return i
	}

	str := fmt.Sprintf("// iogo do\n// %s// %s// %s", s1, s2, s3)
	*pgstr += str

	*varstring += fmt.Sprintf("var Call%s %s = &%sClient{}\n", cliName, cliIface, cliType)

	str = `
	
	iogo.RegisterClient("%s", "proto.%s", Call%s, "")`
	str = fmt.Sprintf(str, unitName, cliName, cliName)
	*initstr += str

	return i + 2
}

func checkRegSrv(unitName string, lines []string, i int, pgstr *string, regstr *string) int {
	if len(lines)-i < 3 {
		return i
	}

	s1 := lines[i]
	s2 := lines[i+1]
	s3 := lines[i+2]

	srvName := getBetweenStr(s1, "func Register", "Server(s *grpc.Server, srv ")
	srvType := getBetweenStr(s1, "Server(s *grpc.Server, srv ", ")")
	fileDesc := getBetweenStr(s2, "s.RegisterService(&", ", srv)")
	if srvName == "" || srvType == "" || fileDesc == "" {
		return i
	}

	str := fmt.Sprintf("// iogo do\n// %s// %s// %s", s1, s2, s3)
	*pgstr += str

	str = `

func Register%sServer(srv %sServer) {
	iogo.RegisterServer("%s", "proto.%s", srv, &_%s_serviceDesc)
}`
	str = fmt.Sprintf(str, srvName, srvName, unitName, srvName, srvName)
	*regstr += str

	return i + 2
}

func createIgFile(filename, vars, inits, regs string) {
	str := `package proto

import "context"

import "iogo"
import "iogo/unit/md"

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

%s
func onStart() {

}

func init() {
	iogo.RegisterFuncOnStart(onStart)%s
}%s
`
	str = fmt.Sprintf(str, vars, inits, regs)
	StrToFile(str, filename)
}

func writeIogoGrpc(unitName, filename string) bool {
	name := strings.Replace(path.Base(filename), ".pb.go", "", -1)
	igFileName := strings.Replace(filename, ".pb.go", ".ig.go", -1)
	name = firstUp(name)

	lines := fileToLine(filename)
	if len(lines) == 0 {
		return true
	}

	pbstr := ""
	varstr := ""
	initstr := ""
	regstr := ""
	i := 0

	for {
		if i >= len(lines) {
			break
		}
		ret := checkRegSrv(unitName, lines, i, &pbstr, &regstr)
		if ret != i {
			i = ret + 1
			continue
		}
		ret = checkNewCli(unitName, lines, i, &pbstr, &varstr, &initstr)
		if ret != i {
			i = ret + 1
			continue
		}

		//str := strings.Replace(lines[i], "(context.Context,", "(Context,", -1)
		//str = strings.Replace(str, "ctx context.Context", "ctx Context", -1)
		//str = strings.Replace(str, "interceptor(ctx,", "interceptor(ctx.(context.Context),", -1)
		pbstr += lines[i]
		i++
	}
	ret := false
	if varstr != "" && initstr != "" && regstr != "" {
		createIgFile(igFileName, varstr, initstr, regstr)
		ret = true
		pbstr = strings.Replace(pbstr, "//iogoimport", "import", 1)
	}
	StrToFile(pbstr, filename)
	return ret
}

func isHttpNet(clusterdir string) bool {
	fn := clusterdir + "/comm/iogo/iogo.go"
	lines := fileToLine(fn)
	for _, line := range lines {
		if strings.Index(line, "/comm/iogo/http") != -1 {
			return true
		}
	}
	return false
}

func setImportGrpc(protodir string) {
	fn := strings.Replace(protodir, "proto", "comm/iogo/iogo.go", -1)
	lines := fileToLine(fn)
	if len(lines) == 0 {
		return
	}
	str := ""
	for _, line := range lines {
		str += line
	}
	str = strings.Replace(str, "/comm/iogo/http", "/comm/iogo/grpc", -1)
	StrToFile(str, fn)
}

func checkServiceDefine(filepath, filename string) (retpath, retname string) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	filestr := ""
	rd := bufio.NewReader(f)
	count := 0
	methods := make(map[string]*[]string)
	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		count++
		if strings.Index(line, "service") > -1 {
			return filepath, filename
		}
		filestr = filestr + line
		if strings.Index(line, "message") != 0 {
			continue
		}
		str := strings.Trim(line, "\n")
		str = strings.Trim(str, "\r")
		msgname := getMsgFullName(str)
		if msg := getMsgName(msgname, "Request"); msg != "" {
			method := methods[msg]
			if method == nil {
				m := make([]string, 2)
				method = &m
				methods[msg] = &m
			}
			(*method)[0] = msg
		}
		if msg := getMsgName(msgname, "Reply"); msg != "" {
			method := methods[msg]
			if method == nil {
				m := make([]string, 2)
				method = &m
				methods[msg] = &m
			}
			(*method)[1] = msg
		}
	}
	name := path.Base(filename)
	service := fmt.Sprintf("\nservice %s {\n", firstUp(strings.Replace(name, ".proto", "", -1)))
	for k, v := range methods {
		req := (*v)[0]
		res := (*v)[1]
		if (req != "") && (res == "") {
			fmt.Printf("Error: The message %sRequest not define Reply in file %s\n", k, filename)
			continue
		}
		if (req == "") && (res != "") {
			fmt.Printf("Error: The message %sReply not define Request in file %s\n", k, filename)
			continue
		}
		service = service + fmt.Sprintf("    rpc %s (%sRequest) returns (%sReply) {\n    }\n", k, req, res)
	}
	service = service + "}\n\n"
	filestr = filestr + service

	retpath = os.TempDir() + "/src"
	if runtime.GOOS == "windows" {
		retpath = strings.Replace(retpath, "\\", "/", -1)
	}
	subdir := strings.Replace(path.Dir(filename), filepath, "", -1)
	filedir := retpath + subdir
	if subdir != "" {
		os.MkdirAll(filedir, os.ModePerm)
	}
	retname = filedir + "/" + name
	StrToFile(filestr, retname)

	return
}

func execProtoc(protodir string, isGrpc bool) {
	if !PathExists(protodir) {
		return
	}
	protoPath := protodir + "/src"
	pbfiles := getProtoFiles(protoPath)
	os.Chdir(protodir)
	for _, filename := range pbfiles {
		outdir := os.TempDir()
		if runtime.GOOS == "windows" {
			outdir = strings.Replace(outdir, "\\", "/", -1)
			filename = strings.Replace(filename, "\\", "/", -1)
		}
		newPath := protoPath
		newName := filename
		newPath, newName = checkServiceDefine(newPath, newName)
		cmd := fmt.Sprintf("protoc --go_out=%s %s --proto_path %s", protodir, newName, newPath)
		if isGrpc {
			cmd = fmt.Sprintf("protoc --go_out=plugins=grpc:%s %s --proto_path %s", protodir, newName, newPath)
		}
		if execCommand(cmd) == false {
			continue
		}
		pbgoName := strings.Replace(path.Base(filename), ".proto", ".pb.go", -1)
		pbgoDir := strings.Replace(path.Dir(filename), protoPath, protodir, -1)
		pbgoName = pbgoDir + "/" + pbgoName
		unitName := strings.Replace(path.Base(filename), ".proto", "", -1)
		dirs := strings.Split(pbgoDir, "/")
		unitName = fmt.Sprintf("%s.%s.proto", dirs[len(dirs)-1], unitName)
		if isGrpc {
			writeIogoGrpc(unitName, pbgoName)
		} else {
			writeIogoHttp(unitName, filename, pbgoName)
		}
	}
	if isGrpc {
		setImportGrpc(protodir)
	} else {
		setImportHttp(protodir)
	}
}
