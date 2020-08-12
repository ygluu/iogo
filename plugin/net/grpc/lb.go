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
	"context"
	"errors"
	"sync"

	"google.golang.org/grpc"

	"iogo/plugin"
)

type hashBalan struct {
	target string
	lb     plugin.LoadBalan
	mu     sync.Mutex
	addrCh chan []grpc.Address
	done   bool
}

func NewBalaner(lb plugin.LoadBalan) grpc.Balancer {
	return &hashBalan{lb: lb}
}

func (self *hashBalan) watchAddrUpdates() error {
	ch := self.lb.Notify()
	<-ch
	addrs := self.lb.Address()
	open := make([]grpc.Address, len(addrs))
	for i, v := range addrs {
		open[i].Addr = v.Addr
		open[i].Metadata = v.Metadata
	}
	if self.done {
		return errors.New("grpc: the client connection is closing")
	}
	select {
	case <-self.addrCh:
	default:
	}
	self.addrCh <- open
	return nil
}

func (self *hashBalan) Start(target string, config grpc.BalancerConfig) error {
	self.mu.Lock()
	defer self.mu.Unlock()
	if self.done {
		return errors.New("grpc: the client connection is closing")
	}
	if self.lb == nil {
		self.target = target
		return nil
	}
	self.addrCh = make(chan []grpc.Address, 1)
	go func() {
		for {
			if err := self.watchAddrUpdates(); err != nil {
				return
			}
		}
	}()
	return nil
}

func (self *hashBalan) Up(addr grpc.Address) func(error) {
	self.lb.Up(addr.Addr, addr.Metadata)
	return func(err error) {
		self.down(addr, err)
	}
}

func (self *hashBalan) down(addr grpc.Address, err error) {
	self.lb.Down(addr.Addr)
}

func (self *hashBalan) Get(ctx context.Context, opts grpc.BalancerGetOptions) (addr grpc.Address, put func(), err error) {
	if self.target != "" {
		addr.Addr = self.target
		return
	}
	addr.Addr = self.lb.Get("")
	if addr.Addr == "" {
		err = errors.New("there is no address available")
	}
	return
}

func (self *hashBalan) Notify() <-chan []grpc.Address {
	return self.addrCh
}

func (self *hashBalan) Close() error {
	self.mu.Lock()
	defer self.mu.Unlock()
	if self.done {
		return errors.New("grpc: the client connection is closing")
	}
	self.done = true
	if self.addrCh != nil {
		close(self.addrCh)
	}
	return nil
}
