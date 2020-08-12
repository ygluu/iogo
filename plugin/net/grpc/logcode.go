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

package grpc // import "iogo/plugin/net/grpc"

import (
	"errors"

	"iogo/log"
	"iogo/plugin"
)

const (
	LC_GRPCS_READY = plugin.LC_GRPC_MODULE + iota
	LC_GRPCS_STOP
	LC_GRPCC_READY
	LC_GRPCC_STOP
	LC_GRPCS_LISTENERR
	LC_GRPCC_CODEERR
	LC_GRPCC_DIALERR
	LC_GRPCC_TLSFILEERR
	LC_GRPCS_TLSFILEERR
	LC_GRPCC_TLS
	LC_GRPCC_CATLS
	LC_GRPCS_TLS
	LC_GRPCS_CATLS
)

func newError(ecode int, estr string) error {
	log.E(ecode, estr)
	return errors.New(estr)
}
