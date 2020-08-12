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

package etcd // import "iogo/plugin/sd/etcd"

import (
	"iogo/plugin"
)

const (
	LC_HTTPS_READY = plugin.LC_ETCD_MODULE + iota
	LC_ETCDS_REGERR
	LC_ETCDS_REGOK
	LC_ETCDS_GETKEYERR
	LC_ETCD_NEWETCDCLIENTERR
	LC_ETCDC_GETERR
	LC_ETCDC_READY
	LC_ETCDC_STOP
	LC_ETCDS_REGING
	LC_ETCDS_STOP
	LC_ETCDC_ONINLINE
)
