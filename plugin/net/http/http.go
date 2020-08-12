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
	"crypto/tls"
	"net"
	"net/http"
)

type NewOption interface {
	Flag() string
	Matadata() interface{}
}

type newOption struct {
	flag  string
	mdata interface{}
}

func (self *newOption) Flag() string {
	return self.flag
}

func (self *newOption) Matadata() interface{} {
	return self.mdata
}

const flag_certfilesvr = "certfilesvr"
const flag_certfilecli = "certfilecli"
const flag_tlscfg = "tlscfg"
const flag_trans = "trans"
const flag_listen = "listen"
const flag_addr = "addr"
const flag_server = "server"
const flag_client = "client"
const flag_default = "default"

func WithDefault() NewOption {
	return &newOption{flag: flag_default, mdata: nil}
}

func WithCertFilesSvr(addr, certFile, keyFile, caFile string) NewOption {
	md := make(map[string]string)
	md["addr"] = addr
	md["cafile"] = caFile
	md["certfile"] = certFile
	md["keyfile"] = keyFile
	return &newOption{flag: flag_certfilesvr, mdata: md}
}

func WithCertFilesCli(certFile, keyFile, caFile, certSvrName string) NewOption {
	md := make(map[string]string)
	md["certsvrname"] = certSvrName
	md["cafile"] = caFile
	md["certfile"] = certFile
	md["keyfile"] = keyFile
	return &newOption{flag: flag_certfilecli, mdata: md}
}

type CfgTls = tls.Config

func WithTlsConfig(cfg *CfgTls) NewOption {
	return &newOption{flag: flag_tlscfg, mdata: cfg}
}

type HttpTransport = http.Transport

func WithTransport(trans *HttpTransport) NewOption {
	return &newOption{flag: flag_trans, mdata: trans}
}

type HttpListen = net.Listener

func WithListen(httpListen *HttpListen) NewOption {
	return &newOption{flag: flag_listen, mdata: httpListen}
}

func WithAddr(addr string) NewOption {
	return &newOption{flag: flag_addr, mdata: addr}
}

type HttpServer = http.Server

func WithHttpServer(httpServer *HttpServer) NewOption {
	return &newOption{flag: flag_server, mdata: httpServer}
}

type HttpClient = http.Client

func WithHttpClient(httpClient *HttpClient) NewOption {
	return &newOption{flag: flag_client, mdata: httpClient}
}
