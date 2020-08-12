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

package http // import "iogo/plugin/net/http"

import (
	"iogo/plugin"
)

const (
	LC_HTTPS_READY = plugin.LC_HTTP_MODULE + iota
	LC_HTTPS_STOP
	LC_HTTPC_READY
	LC_HTTPC_STOP
	LC_HTTPC_READCARTFILEERR
	LC_HTTPC_NOTSTRAT
	LC_HTTPC_REQUESTMARSHALERR
	LC_HTTPC_NEWREQUESTERR
	LC_HTTPC_SVCOUTLINE
	LC_HTTPC_SVCDIALERR
	LC_HTTPC_CALLRESULTUNMARSHALERR
	LC_HTTPC_CALLNOTRESPONSE
	LC_HTTPC_RESULTUNMARSHALERR
	LC_HTTPS_URLERR
	LC_HTTPS_NOTMETHOD
	LC_HTTPC_REQUESTUNMARSHALERR
	LC_HTTPC_INVALIDREQUEST
	LC_HTTPS_READBODYERR
	LC_HTTPS_SERVEHTTPERR
	LC_HTTPS_CALLMETHODERR
	LC_HTTPS_CALLBEFORERR
	LC_HTTPS_CALLAFTERERR
	LC_HTTPS_LISTENERR
	LC_HTTPS_HEADERERR
	LC_HTTPS_METHODNAMEERR
	LC_HTTPS_REQUESTBODYERR
	LC_HTTPS_CALLSVCMETHODERR
	LC_HTTPS_SVCMETHODNOTRETURN
	LC_HTTPC_CODEERR
	LC_HTTPC_CERTFILEERR
	LC_HTTPS_
)

const IOGO_ERROR_NAME = "IOGOERROR"
const IOGO_ERROR_VALUE = "TRUE"
