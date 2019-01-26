package main

import (
	"fmt"
	"os"
	"strings"
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return (err == nil)
}

func StrToFile(s string, name string) {
	file, err := os.Create(name)
	if err != nil {
		fmt.Println("File is exists: ", name)
		return
	}
	defer file.Close()
	file.WriteString(s)
}

func copyFile(dir, filestr, name, newclus, newsvc string) {
	os.MkdirAll(dir, os.ModePerm)
	if newclus != "" {
		filestr = strings.Replace(filestr, "iogo-demo", newclus, -1)
	}
	if newsvc != "" {
		filestr = strings.Replace(filestr, "svc-demo", newsvc, -1)
	}
	StrToFile(filestr, dir+name)
}

func copySvc(clusDir, svc string) {
	copyFile(clusDir+svc+"/", mainFile, "main.go", *clus, svc)
	copyFile(clusDir+svc+"/service/", clientFile, "client.go", *clus, svc)
	copyFile(clusDir+svc+"/service/", serverFile, "server.go", *clus, svc)

}

func createSvc() {
	clusDir := *srcdir + *clus + "/"
	isFirst := false
	if PathExists(clusDir) == false {
		fmt.Printf("create cluster:%s in dir:%s\n", *clus, *srcdir)
		os.MkdirAll(clusDir, os.ModePerm)
		copyFile(clusDir+"config/", confFile, "config.go", *clus, *csvc)
		copyFile(clusDir+"iogo/", iogoFile, "iogo.go", *clus, "")
		copyFile(clusDir+"proto/", protoFile, "hello.proto", "", "")
		copyFile(clusDir+"test-client/", testCliFile, "client.go", *clus, *csvc)
		copyFile(clusDir+"test-client/", testMainFile, "main.go", *clus, *csvc)
		copySvc(clusDir, "svc-hello")
		isFirst = true
	}
	fmt.Printf("create service:%s in dir:%s\n", *csvc, clusDir)
	copySvc(clusDir, *csvc)
	if isFirst {
		os.Chdir(clusDir + "proto/")
		createProto()
	}
}
