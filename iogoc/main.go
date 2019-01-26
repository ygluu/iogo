package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	csvc   = flag.String("csvc", "", "Create a service project under gopath with the value of the project name, see clus parameter.")
	clus   = flag.String("clus", "", "Cluster name. The default is empty, indicating that the service project directory is in the gopath directory.")
	cpro   = flag.String("cpro", "*", "Create golang's protocol source file, which is the name of the protocol file and defaults to all protocol files in the current directory.See pdir and mdir for the source file output directory.")
	pdir   = flag.String("pdir", ".", "Proto(xxx.pb.go) file output directory, default to current directory.")
	mdir   = flag.String("mdir", ".", "Method(xxx.go) file output directory, default to current directory.")
	srcdir = flag.String("srcdir", "", "Source directory. The default is the last directory in the environment GOPATH.")
)

func main() {
	flag.Parse()

	if *srcdir == "" {
		gopaths := strings.Split(os.Getenv("GOPATH"), ";")
		*srcdir = gopaths[len(gopaths)-1]
		*srcdir = strings.Replace(*srcdir, "\\", "/", -1)
		if []byte(*srcdir)[len(*srcdir)-1] != '/' {
			*srcdir = *srcdir + "/"
		}
		*srcdir = *srcdir + "src/"
	}

	//*clus= "iogo-test"
	//*csvc = "iogo-svc"

	if (*csvc == "") && (*cpro == "") {
		fmt.Println("No operation command")
		return
	}

	if *csvc != "" {
		if *clus == "" {
			fmt.Println("press input cluster name")
			return
		}
		createSvc()
		return
	}
	if *cpro == "*" {
		createProto()
	}
}
