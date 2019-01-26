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
 *   readme: https://blog.csdn.net/guestcode/article/details/86655540
 */

package iogo

/*****************************************************************************/
// Service discovery
type DiscServer interface {
	Register(clusterName, serviceName, serviceAddr string) error
	Unregister()
	IsInline() bool
}

type DiscClient interface {
	NewResolver(clusterName, serviceName string) interface{}
	GetTarget() string
	IsInline() bool
}

/*****************************************************************************/
// iogo
func Start() {
	startCli()
	startSvr()
}

func Stop() {
	stopCli()
	stopSvr()
}
