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
	"fmt"
	"strings"
	"time"

	etcd3 "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"golang.org/x/net/context"

	"iogo/log"
	"iogo/plugin"
)

type server struct {
	interval int
	ttl      int
}

func (self *server) Clone() plugin.DiscServer {
	ret := &server{}
	*ret = *self
	return ret
}

func (self *server) doRegister(clusterName string, serviceName string, serviceAddr string, funcGetWeight plugin.FuncGetWeight) {

	svcKey := strings.ToLower(fmt.Sprintf("/%s/services/%s/%s", clusterName, serviceName, serviceAddr))
	getLogHead := func() string {
		return fmt.Sprintf("EtcdServer: ")
	}
	log.I(LC_ETCDS_REGING, getLogHead()+"Key is registing ... key: %s", svcKey)
	defer log.I(LC_ETCDS_STOP, getLogHead()+"Register is stop key: %s", svcKey)

	var client *etcd3.Client = nil
	for {
		if gBaseEtcd.client != nil {
			client = gBaseEtcd.client
			break
		}
		time.Sleep(time.Second / 10)
	}

	tm := time.Second * time.Duration(self.interval)
	ticker := time.NewTicker(tm)

	isFirstRegOk := true
	isFirstRegE := true
	isFirstGetE := true

	checkErr := func(err error) {
		if err != nil {
			if isFirstRegE {
				isFirstRegE = false
				log.E(LC_ETCDS_REGERR, getLogHead()+"Key '%s' is regist error to '%s': %s",
					svcKey, EndpointsToStr(gBaseEtcd.cfg.Endpoints), err.Error())
			}
		} else {
			if isFirstRegOk {
				isFirstRegOk = false
				log.I(LC_ETCDS_REGOK, getLogHead()+"Key '%s' is regist successful to '%s'",
					svcKey, EndpointsToStr(gBaseEtcd.cfg.Endpoints))
			}
		}
	}

	getWeight := func() string {
		if funcGetWeight == nil {
			return "0"
		}
		return fmt.Sprintf("%d", funcGetWeight())
	}

	for {
		resp, _ := client.Grant(context.TODO(), int64(self.ttl))
		_, err := client.Get(context.Background(), svcKey)

		if err == nil {
			_, err := client.Put(context.Background(), svcKey, getWeight(), etcd3.WithLease(resp.ID))
			checkErr(err)
		} else if err == rpctypes.ErrKeyNotFound {
			_, err := client.Put(context.TODO(), svcKey, getWeight(), etcd3.WithLease(resp.ID))
			checkErr(err)
		} else {
			if isFirstGetE {
				isFirstGetE = false
				log.E(LC_ETCDS_GETKEYERR, getLogHead()+"Get key Info: %s", err.Error())
			}
		}

		<-ticker.C
	}
}

func (self *server) Register(clusterName string, serviceName string, serviceAddr string, funcGetWeight plugin.FuncGetWeight) {
	go self.doRegister(clusterName, serviceName, serviceAddr, funcGetWeight)
}

func NewServer(newOpt NewOption, interval int, ttl int) plugin.DiscServer {
	newBase(newOpt)
	if interval <= 0 {
		interval = 20
		if ttl < 20 {
			ttl = 25
		}
	}
	if ttl <= 0 {
		ttl = 25
		if ttl < interval {
			ttl = interval + 5
		}
	}
	if ttl < interval {
		ttl = interval + 5
	}
	ret := &server{
		interval: interval,
		ttl:      ttl,
	}
	return ret
}
