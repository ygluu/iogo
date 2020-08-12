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
	"crypto/tls"
	"net"

	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	mytls "iogo/unit/tls"
)

func GetCredsServerOption(certFile, keyFile, caFile string) (ggrpc.ServerOption, error) {
	if certFile == "" {
		return nil, nil
	}
	cfg, err := mytls.NewTlsServer(certFile, keyFile, caFile)
	if err != nil {
		return nil, err
	}
	creds := credentials.NewTLS(cfg)
	ret := ggrpc.Creds(creds)
	return ret, nil
}

func GetCredsClientOption(certFile, keyFile, caFile, certSvrName string) (ggrpc.DialOption, error) {
	if certFile == "" {
		return nil, nil
	}
	cfg, err := mytls.NewTlsClient(certFile, keyFile, caFile, certSvrName)
	if err != nil {
		return nil, err
	}
	creds := credentials.NewTLS(cfg)
	ret := ggrpc.WithTransportCredentials(creds)
	return ret, nil
}

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

func WithDefault(gsos ...GrpcDailOption) NewOption {
	return &newOption{flag: flag_default, mdata: gsos}
}

func WithCertFilesSvr(addr, certFile, keyFile, caFile string, gsos ...GrpcServerOption) NewOption {
	fs := make(map[string]string)
	fs["addr"] = addr
	fs["cafile"] = caFile
	fs["certfile"] = certFile
	fs["keyfile"] = keyFile
	md := []interface{}{fs, gsos}
	return &newOption{flag: flag_certfilesvr, mdata: md}
}

type GrpcDailOption = ggrpc.DialOption

func WithCertFilesCli(certFile, keyFile, caFile, certSvrName string, gdos ...GrpcDailOption) NewOption {
	fs := make(map[string]string)
	fs["certsvrname"] = certSvrName
	fs["cafile"] = caFile
	fs["certfile"] = certFile
	fs["keyfile"] = keyFile
	md := []interface{}{fs, gdos}
	return &newOption{flag: flag_certfilecli, mdata: md}
}

type CfgTls = tls.Config
type GrpcServerOption = ggrpc.ServerOption

func WithTlsConfig(cfg *CfgTls, gsos ...GrpcServerOption) NewOption {
	md := []interface{}{cfg, gsos}
	return &newOption{flag: flag_tlscfg, mdata: md}
}

type GrpcListen = net.Listener

func WithListen(grpcListen *GrpcListen) NewOption {
	return &newOption{flag: flag_listen, mdata: grpcListen}
}

func WithAddr(addr string, gsos ...GrpcServerOption) NewOption {
	md := []interface{}{addr, gsos}
	return &newOption{flag: flag_addr, mdata: md}
}

type GrpcServer = ggrpc.Server

func WithGrpcServer(svr *GrpcServer, ln *GrpcListen) NewOption {
	md := []interface{}{svr, ln}
	return &newOption{flag: flag_server, mdata: md}
}
