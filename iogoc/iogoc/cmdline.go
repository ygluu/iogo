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
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
)

func getIniValue(filename, section, key string) string {
	if filename == "" {
		return ""
	}
	file, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer file.Close()

	key = strings.ToLower(key)
	section = strings.ToLower(section)
	reader := bufio.NewReader(file)
	var sname string = ""
	for {
		linestr, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		linestr = strings.TrimSpace(linestr)
		if linestr == "" {
			continue
		}
		if linestr[0] == ';' {
			continue
		}
		if linestr[0] == '[' && linestr[len(linestr)-1] == ']' {
			sname = linestr[1 : len(linestr)-1]
		} else if strings.ToLower(sname) == section {
			pair := strings.Split(linestr, "=")
			if len(pair) == 2 {
				kname := strings.TrimSpace(pair[0])
				if strings.ToLower(kname) == key {
					return strings.TrimSpace(pair[1])
				}
			}
		}
	}
	return ""
}

func saveLastWorkDir(workdir string) {
	s := `[workdir]
last=%s
	`
	exe := os.Args[0]
	if runtime.GOOS == "windows" {
		exe = strings.Replace(exe, "\\", "/", -1)
	}
	dir := path.Dir(exe)
	workdir = strings.Replace(workdir, "/src", "", -1)
	StrToFile(fmt.Sprintf(s, workdir), dir+"/cfg.ini")
}

func getLastWorkDir() string {
	exe := os.Args[0]
	if runtime.GOOS == "windows" {
		exe = strings.Replace(exe, "\\", "/", -1)
	}
	dir := path.Dir(exe)
	return getIniValue(dir+"/cfg.ini", "workdir", "last")
}

func listDir(p string) (ret []string) {
	count := 0
	dir, err := ioutil.ReadDir(p)
	if err != nil {
		return
	}
	for _, fi := range dir {
		if fi.IsDir() {
			count++
			ret = append(ret, fi.Name())
			//fmt.Println(count, ":", fi.Name())
		}
	}
	return ret
}

func checkIndex(s string, max int) int {
	index, err := strconv.Atoi(s)
	if err != nil {
		return 0
	} else if index >= 1 && index <= max {
		return index
	} else {
		return 0
	}
}

func checkInput(s string) string {
	s = strings.Trim(s, " ")
	s = strings.Replace(s, "  ", " ", -1)
	s = strings.Replace(s, "  ", " ", -1)
	return s
}

func checkCmd(s string, cmd string) []string {
	arr := strings.Split(s, " ")
	if len(arr) != 2 || cmd != arr[0] {
		return []string{}
	}
	return arr
}

func toService(input *bufio.Scanner, workdir string, cluster string) bool {
	clusterdir := workdir + "/" + cluster
	services := listDir(clusterdir)
	printServices := func() {
		fmt.Println("--------------------------------------------------")
		if len(services) > 0 {
			fmt.Println("Services:")
			for _, v := range services {
				if strings.Index(v, "svc") == 0 {
					fmt.Println(v)
				}
			}
		}
		fmt.Println("")
		fmt.Println("Exit please enter: exit")
		fmt.Println("Go back up one level: up")
		fmt.Println("Make a service directory: mks name")
		fmt.Println("Build iogo source file of http: igc")
		fmt.Println("Build protobuf source file of http : pbc")
		fmt.Println("Build protobuf source file of grpc : pbc -g")
		fmt.Println("")
	}
	printServices()
	for {
		fmt.Printf(clusterdir + ":")
		if !input.Scan() {
			break
		}
		line := input.Text()
		line = checkInput(line)
		if line == "?" {
			printServices()
			continue
		}
		if line == "up" {
			break
		}
		/*if line == "protoc" {
			execProtoc(clusterdir + "/proto")
			continue
		}*/
		if line == "igc" {
			execIogoc(clusterdir + "/proto")
			continue
		}
		if line == "pbc" {
			execProtoc(clusterdir+"/proto", false)
			continue
		}
		if line == "exit" {
			return true
		}
		arr := checkCmd(line, "mks")
		if len(arr) == 2 {
			makes(workdir, cluster, arr[1])
			execProtoc(clusterdir+"/proto", !isHttpNet(clusterdir))
			continue
		}
		arr = checkCmd(line, "pbc")
		if len(arr) >= 1 && len(arr) <= 2 {
			arg := ""
			if len(arr) == 2 {
				arg = arr[1]
			}
			if arg == "" || arg == "-g" {
				execProtoc(clusterdir+"/proto", arg == "-g")
				continue
			}
		}
		fmt.Printf("Invalid cluster command: %s\r\n", line)
	}

	return false
}

func toCluster(input *bufio.Scanner, workdir string) bool {
	clusters := listDir(workdir)
	isCluster := func(name string) bool {
		for _, v := range clusters {
			if v == name {
				return true
			}
		}
		return false
	}
	printClusters := func() {
		fmt.Println("--------------------------------------------------")
		if len(clusters) > 0 {
			fmt.Println("Clusters:")
			for i, v := range clusters {
				fmt.Println(fmt.Sprintf("%d: %s", i+1, v))
			}
		}
		fmt.Println("")
		fmt.Println("Exit please enter: exit")
		fmt.Println("Go back up one level: up")
		fmt.Println("Make a cluster directory: mkc name")
		if len(clusters) > 0 {
			fmt.Println("")
			fmt.Println("Please enter index or name of the cluster")
		}
		fmt.Println("")
	}
	printClusters()
	for {
		fmt.Print(workdir + ":")
		if !input.Scan() {
			break
		}
		line := input.Text()
		line = checkInput(line)
		if line == "?" {
			printClusters()
			continue
		}
		if line == "up" {
			break
		}
		if line == "exit" {
			return true
		}
		cluster := ""
		index := checkIndex(line, len(clusters))
		if index != 0 {
			cluster = clusters[index-1]
		} else if isCluster(line) {
			cluster = line
		}
		arr := checkCmd(line, "mkc")
		if len(arr) == 2 {
			cluster = arr[1]
			if !makec(workdir, cluster) {
				continue
			}
			execProtoc(fmt.Sprintf("%s/%s/proto", workdir, cluster), true)
		}
		if cluster != "" {
			saveLastWorkDir(workdir)
			if toService(input, workdir, cluster) {
				return true
			}
			printClusters()
			continue
		}
		fmt.Printf("Invalid workspace command: %s\r\n", line)
	}
	return false
}

func toWorkDir(input *bufio.Scanner) bool {
	gopaths := strings.Split(os.Getenv("GOPATH"), ";")
	printGopaths := func() {
		fmt.Println("--------------------------------------------------")
		fmt.Println("GoPaths:")
		for i, v := range gopaths {
			s := fmt.Sprintf("%d: %s", i+1, v)
			if runtime.GOOS == "windows" {
				s = strings.Replace(s, "\\", "/", -1)
			}
			fmt.Println(s)
		}
		fmt.Println("")
		fmt.Println("Enter gopath index or gopath full name")
		fmt.Println("")
	}
	workdir := getLastWorkDir()
	if workdir != "" && PathExists(workdir) {
		workdir = workdir + "/src"
		if !PathExists(workdir) {
			os.Mkdir(workdir, os.ModePerm)
		}
		if toCluster(input, workdir) {
			return true
		}
	}
	workdir = ""
	printGopaths()
	for {
		fmt.Print("Please enter gopath:")
		if !input.Scan() {
			break
		}
		line := input.Text()
		line = checkInput(line)
		if line == "up" {
			break
		}
		if line == "exit" {
			return true
		}
		if line == "?" {
			help()
			continue
		}
		var workdir string
		index := checkIndex(line, len(gopaths))
		if index != 0 {
			workdir = gopaths[index-1]
		} else if PathExists(line) {
			workdir = line
		} else {
			fmt.Printf("Invalid command:%s\r\n", line)
			continue
		}
		if runtime.GOOS == "windows" {
			workdir = strings.Replace(workdir, "\\", "/", -1)
		}
		workdir = strings.TrimRight(workdir, "/")
		workdir = workdir + "/src"
		if !PathExists(workdir) {
			os.Mkdir(workdir, os.ModePerm)
		}
		if toCluster(input, workdir) {
			return true
		}
		printGopaths()
	}
	return false
}

func help() {
	fmt.Println("--------------------------------------------------")
	fmt.Println("Iogo version: 1.0")
	fmt.Println("Exit please enter: exit")
	fmt.Println("Go back up one level: up")
	fmt.Println("Make a cluster: mks name")
	fmt.Println("Make a service: mkc name")
	fmt.Println("Build iogo source file of http: igc")
	fmt.Println("Build protobuf source file of http : pbc")
	fmt.Println("Build protobuf source file of grpc : pbc -g")
	fmt.Println("--------------------------------------------------")
}

func Cmdline() {
	fmt.Println("--------------------------------------------------")
	fmt.Println("Hello, I'm iogoc!")
	help()
	fmt.Println("Let's get to working")
	input := bufio.NewScanner(os.Stdin)
	for {
		if toWorkDir(input) {
			break
		}
		fmt.Println("Exit please enter: exit")
	}
}
