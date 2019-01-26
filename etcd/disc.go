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

// service discovery plugin: etcd

package etcd

import (
	"errors"
	"fmt"
	"log"
	"time"

	etcd3 "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"golang.org/x/net/context"
	"google.golang.org/grpc/naming"

	"iogo"
)

/*****************************************************************************/
// server

type register struct {
	owner        *server
	client       *etcd3.Client
	serviceName  string
	serviceKey   string
	serviceValue string
	stopSignal   chan bool
	isInline     bool
}

func (self *register) register() error {
	var err error
	self.isInline = false
	isFirst := true

	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(self.owner.interval))
		for {
			if self.client == nil {
				self.client, err = etcd3.New(*self.owner.conf)
				if err != nil {
					self.client = nil
					if isFirst {
						isFirst = false
						log.Printf("iogo-etcd: create etcd3 client failed: %v\n", err)
					}
					continue
				}
			}
			resp, _ := self.client.Grant(context.TODO(), int64(self.owner.ttl))
			_, err := self.client.Get(context.Background(), self.serviceKey)
			if err != nil {
				if err == rpctypes.ErrKeyNotFound {
					if _, err := self.client.Put(context.TODO(), self.serviceKey, self.serviceValue, etcd3.WithLease(resp.ID)); err != nil {
						log.Printf("iogo-etcd: set service '%s' with ttl to etcd3 failed: %s", self.serviceName, err.Error())
					} else {
						if self.isInline == false {
							log.Printf("iogo-etcd: service register: name:%s addr: %s", self.serviceName, self.serviceValue)
						}
						self.isInline = true
					}
				} else {
					log.Printf("iogo-etcd: service '%s' connect to etcd3 failed: %s", self.serviceName, err.Error())
				}
			} else {
				if _, err := self.client.Put(context.Background(), self.serviceKey, self.serviceValue, etcd3.WithLease(resp.ID)); err != nil {
					log.Printf("iogo-etcd: refresh service '%s' with ttl to etcd3 failed: %s", self.serviceName, err.Error())
				} else {
					if self.isInline == false {
						log.Printf("iogo-etcd: service register: name:%s addr: %s", self.serviceName, self.serviceValue)
					}
					self.isInline = true
				}
			}
			select {
			case <-self.stopSignal:
				return
			case <-ticker.C:
			}
		}
	}()

	return nil
}

func (self *register) Unregister() error {
	self.stopSignal <- true
	var err error
	if _, err := self.client.Delete(context.Background(), self.serviceKey); err != nil {
		log.Printf("iogo-etcd: unregister '%s' failed: %s", self.serviceKey, err.Error())
	} else {
		log.Printf("iogo-etcd: unregister '%s' ok.", self.serviceKey)
	}
	return err
}

type server struct {
	interval int
	ttl      int
	conf     *etcd3.Config
	regs     []*register
}

func (self *server) Register(clusterName, serviceName, serviceAddr string) error {
	ret := &register{
		owner:        self,
		serviceName:  serviceName,
		serviceValue: serviceAddr,
		serviceKey:   fmt.Sprintf("/%s/services/%s/%s", clusterName, serviceName, serviceAddr),
		stopSignal:   make(chan bool, 1),
	}
	self.regs = append(self.regs, ret)
	return ret.register()
}

func (self *server) Unregister() {
	for _, reg := range self.regs {
		reg.Unregister()
	}
}

func (self *server) IsInline() bool {
	for _, r := range self.regs {
		if r.isInline {
			return true
		}
	}
	return false
}

func (self *server) getTarget() string {
	ret := ""
	for _, s := range self.conf.Endpoints {
		if ret == "" {
			ret = s
		} else {
			ret = ret + "," + s
		}
	}
	return ret
}

func NewServer(interval int, ttl int, conf *etcd3.Config) iogo.DiscServer {
	ret := &server{
		conf:     conf,
		interval: interval,
		ttl:      ttl,
	}
	return ret
}

/*****************************************************************************/
// client

type client struct {
	conf     *etcd3.Config
	watchers []*watcher
	isInline bool
}

func (self *client) NewResolver(clusterName, serviceName string) interface{} {
	return &resolver{clusterName: clusterName, serviceName: serviceName, owner: self}
}

func (self *client) GetTarget() string {
	ret := ""
	for _, s := range self.conf.Endpoints {
		if ret == "" {
			ret = s
		} else {
			ret = ret + "," + s
		}
	}
	return ret
}

func (self *client) IsInline() bool {
	if len(self.watchers) == 0 {
		return true
	}
	return self.isInline
}

func NewClient(conf *etcd3.Config) iogo.DiscClient {
	return &client{
		conf: conf,
	}
}

type resolver struct {
	owner       *client
	clusterName string
	serviceName string
}

func (self *resolver) Resolve(target string) (naming.Watcher, error) {
	if self.serviceName == "" {
		return nil, errors.New("iogo-etcd: no service name provided")
	}
	client, err := etcd3.New(*self.owner.conf)
	if err != nil {
		client = nil
		log.Printf("iogo-etcd: Resolve: creat etcd3 client failed: %s\n", err.Error())
	} else {
		self.owner.isInline = true
	}
	ret := &watcher{re: self, clusterName: self.clusterName, owner: self.owner, client: client}
	self.owner.watchers = append(self.owner.watchers, ret)
	return ret, nil
}

type watcher struct {
	owner         *client
	clusterName   string
	re            *resolver
	client        *etcd3.Client
	isInitialized bool
}

func (self *watcher) Close() {
}

func (self *watcher) extractAddrs(resp *etcd3.GetResponse) []string {
	addrs := []string{}

	if resp == nil || resp.Kvs == nil {
		return addrs
	}

	for i := range resp.Kvs {
		if v := resp.Kvs[i].Value; v != nil {
			addrs = append(addrs, string(v))
		}
	}

	return addrs
}

func (self *watcher) Next() ([]*naming.Update, error) {
	prefix := fmt.Sprintf("/%s/services/%s/", self.clusterName, self.re.serviceName)

	if !self.isInitialized {
		if self.client == nil {
			client, err := etcd3.New(*self.owner.conf)
			if err != nil {
				self.client = nil
				return nil, fmt.Errorf("iogo-etcd: Next: creat etcd3 client failed: %s", err.Error())
			}
			self.owner.isInline = true
			self.client = client
		}
		resp, err := self.client.Get(context.Background(), prefix, etcd3.WithPrefix())
		if err == nil {
			self.isInitialized = true
			addrs := self.extractAddrs(resp)
			if l := len(addrs); l != 0 {
				log.Printf("iogo-etcd: watcher services name:%s count:%d", self.re.serviceName, l)
				updates := make([]*naming.Update, l)
				for i := range addrs {
					updates[i] = &naming.Update{Op: naming.Add, Addr: addrs[i]}
				}
				return updates, nil
			}
		} else {
			return nil, nil
		}
	}

	rch := self.client.Watch(context.Background(), prefix, etcd3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				//log.Printf("iogo-etcd: watcher services name:%s addr:%s", self.re.serviceName, string(ev.Kv.Value))
				return []*naming.Update{{Op: naming.Add, Addr: string(ev.Kv.Value)}}, nil
			case mvccpb.DELETE:
				return []*naming.Update{{Op: naming.Delete, Addr: string(ev.Kv.Value)}}, nil
			}
		}
	}
	return nil, nil
}
