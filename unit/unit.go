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

package unit // import "iogo/unit"

import (
	"net"
	"reflect"
)

func GetServiceIP(subnetMask string) string {
	if subnetMask == "" {
		return ""
	}
	_, IPNet, err := net.ParseCIDR(subnetMask)
	if err != nil {
		return ""
	}

	getIP := func(addr net.Addr) net.IP {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		if ip == nil || ip.IsLoopback() {
			return nil
		}
		ip = ip.To4()
		if ip == nil {
			return nil
		}
		return ip
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return ""
		}
		for _, addr := range addrs {
			ip := getIP(addr)
			if ip == nil {
				continue
			}
			if !IPNet.Contains(ip) {
				continue
			}
			return ip.String()
		}
	}

	return ""
}

func IsStruct(v interface{}) bool {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Struct {
		return true
	}
	if t.Kind() != reflect.Ptr {
		return false
	}
	return (t.Elem().Kind() == reflect.Struct)
}

func IsBytePoint(v interface{}) bool {
	t := reflect.TypeOf(v)
	return t.String() == "*[]byte"
}

func IsByteDPoint(v interface{}) bool {
	t := reflect.TypeOf(v)
	return t.String() == "**[]byte"
}

func IsStructPoint(v interface{}) bool {
	return IsStructPointByType(reflect.TypeOf(v))
}

func IsStructPointByType(t reflect.Type) bool {
	if t.Kind() != reflect.Ptr {
		return false
	}
	return (t.Elem().Kind() == reflect.Struct)
}
