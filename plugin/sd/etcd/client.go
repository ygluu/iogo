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
	"strconv"
	"strings"
	"time"

	etcd3 "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"golang.org/x/net/context"

	"iogo/log"
	"iogo/plugin"
)

type client struct {
	newOpt NewOption
}

func (self *client) doWatch(clusterName string, serviceName string, lb plugin.LoadBalan) {
	getLogHead := func() string {
		return fmt.Sprintf("Etcdclient(%s/%s): ", clusterName, serviceName)
	}
	prefix := strings.ToLower(fmt.Sprintf("/%s/services/%s/", clusterName, serviceName))
	log.I(LC_ETCDC_READY, "EtcdClient: Watcher '%s' is starting ...", prefix)

	var client *etcd3.Client = nil
	for {
		if gBaseEtcd.client != nil {
			client = gBaseEtcd.client
			break
		}
		time.Sleep(time.Second / 10)
	}

	keyToAddr := func(key string) string {
		arr := strings.Split(key, "/")
		return arr[len(arr)-1]
	}
	isFirstOk := true
	postMsg := func(k string, v string, isAdd bool) {
		addr := keyToAddr(k)
		weight, err := strconv.Atoi(v)
		if err != nil {
			weight = 1
		}
		if isAdd {
			lb.Add(addr, weight)
			if isFirstOk {
				isFirstOk = false
				log.I(LC_ETCDC_READY, "EtcdClient: Watcher is ready key: %s", prefix)
			}
			//log.D(LC_ETCDC_ONINLINE, "EtcdClient: Service inline addr: %s", addr)
		} else {
			lb.Del(addr)
		}
	}

	isFirstE := true
	for {
		resp, err := client.Get(context.Background(), prefix, etcd3.WithPrefix())
		if err != nil {
			if isFirstE {
				isFirstE = false
				log.E(LC_ETCDC_GETERR, getLogHead()+"%+v", err.Error())
			}
			time.Sleep(1 * time.Second)
			continue
		}

		if resp != nil && resp.Kvs != nil {
			for _, kv := range resp.Kvs {
				postMsg(string(kv.Key), string(kv.Value), true)
			}
		}

		rch := client.Watch(context.Background(), prefix, etcd3.WithPrefix())
		for wresp := range rch {
			if wresp.Events == nil {
				return
			}
			for _, ev := range wresp.Events {
				switch ev.Type {
				case mvccpb.PUT:
					postMsg(string(ev.Kv.Key), string(ev.Kv.Value), true)
				case mvccpb.DELETE:
					postMsg(string(ev.Kv.Key), "0", false)
				}
			}
		}
	}
}

func (self *client) Clone() plugin.DiscClient {
	ret := &client{
		newOpt: self.newOpt,
	}
	return ret
}

func (self *client) Watcher(clusterName string, serviceName string, loadBalan plugin.LoadBalan) {
	go self.doWatch(clusterName, serviceName, loadBalan)
}

func NewClient(newOpt NewOption) plugin.DiscClient {
	newBase(newOpt)
	ret := &client{
		newOpt: newOpt,
	}
	return ret
}
