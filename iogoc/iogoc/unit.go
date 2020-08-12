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
	"bytes"
	"fmt"
	"os"
	"os/exec"
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

func execCommand(scmd string) bool {
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
