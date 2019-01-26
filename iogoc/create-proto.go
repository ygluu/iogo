package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func FindSvcs(name string) []string {
	ret := make([]string, 0)

	f, err := os.Open(name)
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
		str := strings.TrimRight(strings.Trim(line, "\n"), "\n")
		str = strings.TrimRight(strings.Trim(str, "\r"), "\r")
		if strings.Index(str, "service") > -1 {
			ss := strings.Split(str, " ")
			if (len(ss) != 3) || (ss[0] != "service") || (ss[2] != "{") {
				continue
			}
			ret = append(ret, ss[1])
		}

	}

	return ret
}

func getMethodName(s string) string {
	cs := ([]byte)(s)
	ret := ""
	isTrue := false
	for _, c := range cs {
		if string(c) == "\"" {
			if isTrue {
				return ret
			} else {
				isTrue = true
				continue
			}
		}
		if isTrue {
			ret = ret + string(c)
		}
	}
	return ret
}

func ReplaceReturnInfo(name string) {
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}

	rd := bufio.NewReader(f)
	method := ""
	mline := 0
	count := 0
	filestr := ""
	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		count++
		str := strings.TrimRight(strings.Trim(line, "\n"), "\n")
		str = strings.TrimRight(strings.Trim(str, "\r"), "\r")

		if strings.Index(str, "err := grpc.Invoke(ctx, \"") > -1 {
			method = getMethodName(str)
			if method != "" {
				mline = count
			}
			filestr = filestr + line
			continue
		}
		if (strings.Index(str, "		return nil, err") > -1) && (count-mline == 2) {
			filestr = filestr + "		return nil, fmt.Errorf(\"call " + method + ": %s\", err.Error())\n"
		} else {
			filestr = filestr + line
		}
	}
	f.Close()

	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error: FileToStr.Open file", name, ":", err)
		return
	}
	defer file.Close()

	n, err := file.WriteString(filestr)
	if (err != nil) || (n != len(filestr)) {
		fmt.Println("Error: FileToStr.WriteString file", name, ":", err)
		return
	}
}

func ReplaceInFile(name string, sold []string, snew []string) bool {
	datas, err := ioutil.ReadFile(name)
	if err != nil {
		fmt.Println("Error: FileToStr.ReadFile file", name, ":", err)
		return false
	}

	s := string(datas)
	for i := 0; i < len(sold); i++ {
		s = strings.Replace(s, sold[i], snew[i], -1)
	}

	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error: FileToStr.Open file", name, ":", err)
		return false
	}
	defer file.Close()

	n, err := file.WriteString(s)
	if (err != nil) || (n != len(s)) {
		fmt.Println("Error: FileToStr.WriteString file", name, ":", err)
		return false
	}

	return true
}

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
	datas := []byte(s)
	datas[0] = ([]byte(strings.ToUpper(string(datas[0]))))[0]
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

func getFileName(name string) string {
	name = path.Base(name)
	ext := path.Ext(name)
	return strings.TrimSuffix(name, ext)
}

var (
	ecount = 0
	wcount = 0
)

func execCommand(scmd string) bool {
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

func createSvcDef(name string) string {
	ret := "."

	f, err := os.Open(name)
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
			fmt.Printf("Service has been defined: %s\n", name)
			return ret
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
	service := fmt.Sprintf("\nservice %s {\n", getFileName(name))
	for k, v := range methods {
		req := (*v)[0]
		res := (*v)[1]
		if (req != "") && (res == "") {
			ecount++
			fmt.Printf("Error: The message %sRequest not define Reply in file %s\n", k, name)
			continue
		}
		if (req == "") && (res != "") {
			wcount++
			fmt.Printf("Warning: The message %sReply not define Request in file %s\n", k, name)
			continue
		}
		service = service + fmt.Sprintf("    rpc %s (%sRequest) returns (%sReply) {\n    }\n", k, req, res)
	}
	service = service + "}\n\n"
	filestr = filestr + service

	ret = strings.Replace(os.TempDir(), "\\", "/", -1)
	tmpfile := ret + "/" + name
	os.Remove(tmpfile)
	StrToFile(filestr, tmpfile)

	return ret
}

func fileExists(name string) bool {
	file, err := os.Open(name)
	if err != nil {
		return false
	}
	defer file.Close()
	return true
}

func createProto() {
	fmt.Println("The proto file is creating...")

	files := []string{}
	if *cpro == "*" {
		fs, _ := filepath.Glob("*")
		for _, name := range fs {
			if path.Ext(name) == ".proto" {
				files = append(files, name)
			}
		}
	} else {
		files = append(files, *cpro)
	}
	for _, name := range files {
		if fileExists(name) {
			fmt.Println("\nProcess file: ", name)

			// exec protoc before
			outdir := createSvcDef(name)

			// exec protoc
			cmd := fmt.Sprintf("%s --go_out=plugins=grpc:%s %s --proto_path %s", "protoc", *pdir, name, outdir)
			if execCommand(cmd) == false {
				continue
			}

			// exec protoc after
			olds := make([]string, 0)
			news := make([]string, 0)
			/*svcs := FindSvcs(name)
			for _, s := range svcs {
				olds = append(olds, fmt.Sprintf("func New%sClient(cc *grpc.ClientConn) %sClient {", s, s))
				news = append(news, fmt.Sprintf("func New%sClient(cc interface{}) %sClient {", s, s))
				olds = append(olds, fmt.Sprintf("	return &%sClient{cc}", strings.ToLower(s)))
				news = append(news, fmt.Sprintf("	return &%sClient{cc.(*grpc.ClientConn)}", strings.ToLower(s)))

				olds = append(olds, fmt.Sprintf("func Register%sServer(s *grpc.Server, srv %sServer) {", s, s))
				news = append(news, fmt.Sprintf("func Register%sServer(s interface{}, srv %sServer) {", s, s))
				olds = append(olds, fmt.Sprintf("	s.RegisterService(&_%s_serviceDesc, srv)", s))
				news = append(news, fmt.Sprintf("	s.(*grpc.Server).RegisterService(&_%s_serviceDesc, srv)", s))
			}*/
			olds = append(olds, "fileDescriptor0")
			news = append(news, "fileDesc"+firstUp(getFileName(name)))
			ReplaceInFile(getFileName(name)+".pb.go", olds, news)
			ReplaceReturnInfo(getFileName(name) + ".pb.go")
		}
	}
	fmt.Printf("\nThe proto file create results: error-request:%d, warning-reply:%d\n", ecount, wcount)
}

/*
func NewAccountClient(cc interface{}) AccountClient {
	return &accountClient{cc.(*grpc.ClientConn)}
}
func RegisterAccountServer(s interface{}, srv AccountServer) {
	s.(*grpc.Server).RegisterService(&_Session_serviceDesc, srv)
}
*/
