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

package log // import "iogo/log"

import (
	"iogo/plugin"
	"iogo/plugin/log"
)

var gLogger plugin.Logger = log.NewLogger()

func I(code int, format string, v ...interface{}) {
	gLogger.I(code, format, v...)
}

func W(code int, format string, v ...interface{}) {
	gLogger.W(code, format, v...)
}

func E(code int, format string, v ...interface{}) {
	gLogger.E(code, format, v...)
}

func D(code int, format string, v ...interface{}) {
	gLogger.D(code, format, v...)
}

func T(code int, format string, v ...interface{}) {
	gLogger.T(code, format, v...)
}

func Get() plugin.Logger {
	return gLogger
}

func Set(logger plugin.Logger) {
	log.Copy(gLogger, logger)
	gLogger = logger
}
