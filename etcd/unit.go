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

package etcd

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"strings"

	//"time"

	etcd3 "github.com/coreos/etcd/clientv3"
)

func NewConf(target, caFile, certFile, keyFile string) (*etcd3.Config, error) {
	var tlsCfg *tls.Config = nil
	if (caFile != "") && (certFile != "") && (keyFile != "") {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, err
		}

		caData, err := ioutil.ReadFile(caFile)
		if err != nil {
			return nil, err
		}

		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(caData)

		tlsCfg = &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      pool,
		}
	}
	return &etcd3.Config{
		Endpoints: strings.Split(target, ","),
		//DialTimeout: 5 * time.Second,
		TLS: tlsCfg,
	}, nil
}
