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
	"crypto/tls"
	"reflect"
	"strings"
	"time"

	etcd3 "github.com/coreos/etcd/clientv3"

	"iogo/log"
	mytls "iogo/unit/tls"
)

type CfgEtcd = etcd3.Config
type CfgTls = tls.Config

/*
func NewCfgTls(certFile, keyFile, caFile string) (*CfgTls, error) {
	if (caFile == "") || (certFile == "") || (keyFile == "") {
		return nil, errors.New("iogo.etcd.NewSTLConf: file name invalid")
	}
	var tlsCfg *tls.Config = nil
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
	return tlsCfg, nil
}*/

func EndpointsToStr(endpoints []string) string {
	ret := ""
	for _, v := range endpoints {
		if ret == "" {
			ret = v
		} else {
			ret = ret + ";" + v
		}
	}
	return ret
}

func NewCfgEtcd(endpoints string) *CfgEtcd {
	ret := &etcd3.Config{}
	ret.Endpoints = strings.Split(endpoints, ";")
	return ret
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

const flag_cfg = "cfg"
const flag_certfile = "certfile"
const flag_endp = "endp"

func WithCertFile(endpoints, certSvrName, certFile, keyFile, caFile string) NewOption {
	if certFile == "" || (keyFile == "" && caFile != "") || (keyFile != "" && caFile == "") {
		panic("WithCredentialFilesOfServer: cert file error")
	}
	md := make(map[string]string)
	md["endpoints"] = endpoints
	md["certsvrname"] = certSvrName
	md["cafile"] = caFile
	md["certfile"] = certFile
	md["keyfile"] = keyFile
	return &newOption{flag: flag_certfile, mdata: md}
}

func WithConfig(cfg *etcd3.Config) NewOption {
	return &newOption{flag: flag_cfg, mdata: cfg}
}

func WithEndpoints(endpoints string) NewOption {
	return &newOption{flag: flag_endp, mdata: endpoints}
}

type baseEtcd struct {
	cfg    *CfgEtcd
	client *etcd3.Client
}

func (self *baseEtcd) Cfg() *CfgEtcd {
	return self.cfg
}

func (self *baseEtcd) doStart(cfg *etcd3.Config) {
	go func() {
		if cfg.DialTimeout == 0 {
			cfg.DialTimeout = 6 * time.Second
		}
		for {
			cli, err := etcd3.New(*cfg)
			if err != nil {
				log.E(LC_ETCD_NEWETCDCLIENTERR, "BaseEtcd: %s", err.Error())
				time.Sleep(10 * time.Second)
				continue
			}
			self.client = cli
			break
		}
		cfg = nil
	}()
}

var gBaseEtcd *baseEtcd = nil

func newBase(newOpt NewOption) {
	if gBaseEtcd != nil {
		return
	}
	if newOpt == nil {
		panic("iogo.Etcd: Configuration parameter cannot be empty")
	}
	var cfg *CfgEtcd
	if newOpt.Flag() == flag_cfg {
		cfg = newOpt.Matadata().(*CfgEtcd)
	} else if newOpt.Flag() == flag_endp {
		cfg = NewCfgEtcd(newOpt.Matadata().(string))
	} else if newOpt.Flag() == flag_certfile {
		md := newOpt.Matadata().(map[string]string)
		certSvrName := md["certsvrname"]
		endpoints := md["endpoints"]
		caFile := md["cafile"]
		certFile := md["certfile"]
		keyFile := md["keyfile"]
		tls, err := mytls.NewTlsClient(certFile, keyFile, caFile, certSvrName)
		if err != nil {
			panic("NewHttpServer: cert file error:" + err.Error())
		}
		cfg = NewCfgEtcd(endpoints)
		cfg.TLS = tls
	} else {
		t := reflect.TypeOf(newOpt.Matadata())
		panic("NewHttpServer: NewOption error: " + t.String())
	}
	if len(cfg.Endpoints) == 0 {
		panic("iogo.Etcd: Endpoints is null")
	}
	gBaseEtcd = &baseEtcd{client: nil, cfg: cfg}
	gBaseEtcd.doStart(cfg)
}
