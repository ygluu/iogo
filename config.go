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

package iogo // import "iogo"

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"iogo/unit"
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

type cfgReader struct {
	iniFileName string
}

func (self *cfgReader) getValue(key, getDef, iniSec, flagDef, flagInfo string) string {
	ret := ""
	flag.StringVar(&ret, key, flagDef, flagInfo)
	flag.Parse()
	if ret == "" {
		ret = getIniValue(self.iniFileName, iniSec, key)
	}
	if ret == "" {
		ret = getDef
	}
	return ret
}

func (self *cfgReader) GetClusterName(def string) string {
	ret := self.getValue("cluster_name", def, "Cluster", "", "name of cluster default: empty string")
	return ret
}

func (self *cfgReader) GetServiceName(def string) string {
	ret := self.getValue("service_name", def, "Cluster", "", "name of service default: empty string")
	return ret
}

func (self *cfgReader) GetServiceAddr(defAddr string, defPort string, defSubnetMask string) string {
	addr := self.getValue("service_addr", defAddr, "Cluster", "", "addr of service default: empty string")
	if addr != "" {
		return addr
	}
	port := self.getValue("service_port", defPort, "Cluster", "", "name of service port default: empty string")
	if port == "" {
		return ""
	}
	mask := self.getValue("subnet_mask", defSubnetMask, "Cluster", "", "mask of subnet default: empty string")
	if port == "" {
		return ""
	}
	ip := unit.GetServiceIP(mask)
	if ip == "" {
		return ""
	}
	return fmt.Sprintf("%s:%s", ip, port)
}

func (self *cfgReader) GetSdCaFileName(def string) string {
	ret := self.getValue("sd ca_file", def, "ServiceDiscovery", "", "name of ca file default: empty string")
	return ret
}

func (self *cfgReader) GetSdKeyFileName(def string) string {
	ret := self.getValue("sdkey_file", def, "ServiceDiscovery", "", "name of key file default: empty string")
	return ret
}

func (self *cfgReader) GetSdCertFileName(def string) string {
	ret := self.getValue("sdcert_file", def, "ServiceDiscovery", "", "name of cert file default: empty string")
	return ret
}

func (self *cfgReader) GetSdEndpoints(def string) string {
	ret := self.getValue("sdendpoints", def, "ServiceDiscovery", "", "endpoints of service discovery default: empty string")
	return ret
}

func (self *cfgReader) GetCaFileName(def string) string {
	ret := self.getValue("ca_file", def, "net", "", "name of ca file default: empty string")
	return ret
}

func (self *cfgReader) GetServerKeyFileName(def string) string {
	ret := self.getValue("svr_key_file", def, "net", "", "name of server key file default: empty string")
	return ret
}

func (self *cfgReader) GetServerCertFileName(def string) string {
	ret := self.getValue("svr_cert_file", def, "net", "", "name of server cert file default: empty string")
	return ret
}

func (self *cfgReader) GetClientKeyFileName(def string) string {
	ret := self.getValue("cli_key_file", def, "net", "", "name of client key file default: empty string")
	return ret
}

func (self *cfgReader) GetClientCertFileName(def string) string {
	ret := self.getValue("cli_cert_file", def, "net", "", "name of client cert file default: empty string")
	return ret
}

func (self *cfgReader) GetCertServerName(def string) string {
	ret := self.getValue("cert_svr_name", def, "net", "", "name of cert server default: empty string")
	return ret
}

func getFileName() string {
	execPath, err := exec.LookPath(os.Args[0])
	if err != nil {
		return ""
	}
	return execPath + "/config/cfg.ini"
}

var CfgReader *cfgReader = &cfgReader{iniFileName: getFileName()}
