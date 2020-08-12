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

package log

import (
	"fmt"
	glog "log"
	"math"
	"reflect"
	"time"

	"iogo/plugin"
)

type logMsg struct {
	time   time.Time
	typ    string
	code   int
	format string
	args   []interface{}
}

func (self *logMsg) getHead(key string) string {
	const logHeadStr1 = "%s %s %d%s | "

	if key != "" {
		key = " | " + key
	}

	timeToStr := func(time *time.Time) string {
		um := int(math.Ceil(float64(time.Nanosecond() / 1000)))
		return fmt.Sprintf("%d/%0.2d/%0.2d %0.2d:%0.2d:%0.2d.%d",
			time.Year(), time.Month(), time.Day(),
			time.Hour(), time.Minute(), time.Second(), um)
	}

	return fmt.Sprintf(
		logHeadStr1,
		timeToStr(&self.time),
		self.typ,
		self.code,
		key,
	)
}

type logger struct {
	chanMsg chan *logMsg
	key     string
	isStop  bool
}

func (self *logger) thread() {
	for {
		msg := <-self.chanMsg
		if msg == nil {
			self.isStop = true
			break
		}
		glog.Printf(msg.getHead(self.key)+msg.format+"\n", msg.args...)
	}
}

func (self *logger) postMsg(typ string, code int, format string, v ...interface{}) {
	checkPrt := func(arg interface{}) interface{} {
		t := reflect.TypeOf(arg)
		if t.Kind() != reflect.Ptr {
			return arg
		}
		newV := reflect.New(t.Elem()).Elem()
		newV.Set(reflect.ValueOf(arg).Elem())
		return newV.Addr().Interface()
	}

	var args []interface{}
	for _, arg := range v {
		if arg == nil {
			args = append(args, "nil")
			continue
		}
		args = append(args, checkPrt(arg))
	}

	msg := &logMsg{
		code:   code,
		typ:    typ,
		time:   time.Now(),
		format: format,
		args:   args,
	}
	self.chanMsg <- msg
}

func (self *logger) I(code int, format string, v ...interface{}) {
	self.postMsg("I", code, format, v...)
}

func (self *logger) W(code int, format string, v ...interface{}) {
	self.postMsg("W", code, format, v...)
}

func (self *logger) E(code int, format string, v ...interface{}) {
	self.postMsg("E", code, format, v...)
}

func (self *logger) D(code int, format string, v ...interface{}) {
	self.postMsg("D", code, format, v...)
}

func (self *logger) T(code int, format string, v ...interface{}) {
	self.postMsg("T", code, format, v...)
}

func (self *logger) SetKey(key string) {
	self.key = key
	glog.SetFlags(0) //log.Ldate | log.Ltime | log.Lmicroseconds
	go self.thread()
}

func Copy(src, dest plugin.Logger) {
	t := reflect.TypeOf(src)
	if t.String() == "*log.logger" {
		tmp := reflect.ValueOf(src).Interface().(*logger)
		tmp.chanMsg <- nil
		for {
			msg := <-tmp.chanMsg
			if msg == nil {
				break
			}
			if msg.typ == "E" {
				dest.E(msg.code, msg.format, msg.args...)
			} else if msg.typ == "W" {
				dest.W(msg.code, msg.format, msg.args...)
			} else if msg.typ == "D" {
				dest.D(msg.code, msg.format, msg.args...)
			} else if msg.typ == "T" {
				dest.T(msg.code, msg.format, msg.args...)
			} else if msg.typ == "I" {
				dest.I(msg.code, msg.format, msg.args...)
			}
		}
		dest.SetKey(tmp.key)
	}
}

func NewLogger() plugin.Logger {
	ret := &logger{
		isStop:  false,
		chanMsg: make(chan *logMsg, 100000),
	}
	ret.I(LC_LOG_READY, "Logger: Logger is ready")
	return ret
}
