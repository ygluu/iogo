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
	"fmt"
	"time"

	"iogo/log"
)

type runtime struct {
	reqCount  int64
	respCount int64
	lastTime  time.Time
	IsPrint   bool
}

func (self *runtime) threadPrint() {
	self.lastTime = time.Now()
	start := time.Now()
	oldReq := self.reqCount
	for {
		time.Sleep(2 * time.Second)
		if !self.IsPrint {
			continue
		}
		newReq, resp := gNetServer.GetConcurrentCount()
		end := time.Now()
		count := newReq - oldReq
		if newReq != oldReq {
			self.lastTime = time.Now()
		}
		concur := time.Duration(0)
		if count != 0 {
			concur = time.Duration(1000*1000*1000) / (end.Sub(start) / time.Duration(count))
		}
		log.I(_LC_CORE_RUNINFO, "IogoCore:  Request:%d  Response:%d  Concurrent(s):%d  LastActive:%s",
			newReq, resp, concur, self.lastTime.String())
		oldReq = newReq
		start = end
	}
}

func (self *runtime) GetClusterName() string {
	return gClusterName
}

func (self *runtime) GetServiceAddr() string {
	return gServiceAddr
}

func (self *runtime) GetServiceKey() string {
	if gServiceAddr == "" {
		return "IogoClient/Only"
	}
	return fmt.Sprintf("%s/%s", gClusterName, gServiceAddr)
}

var Runtime runtime = runtime{
	reqCount:  0,
	respCount: 0,
	IsPrint:   false,
}

func startRuntime() {
	if gNetServer != nil {
		go Runtime.threadPrint()
	}
}
