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
	"golang.org/x/net/context"

	etcd3 "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
)

/*****************************************************************************/
// lock
type Lock struct {
	s *concurrency.Session
	m *concurrency.Mutex
}

func (self *Lock) Lock() error {
	if err := self.m.Lock(context.TODO()); err != nil {
		return err
	}
	return nil
}

func (self *Lock) Unlock() error {
	err := self.m.Unlock(context.TODO())
	if err != nil {
		return err
	}
	return nil
}

func (self *Lock) Close() {
	self.Unlock()
	self.s.Close()
}

func NewLock(key string, conf *etcd3.Config) (*Lock, error) {
	client, err := etcd3.New(*conf)
	if err != nil {
		return nil, err
	}
	lock := &Lock{}
	lock.s, err = concurrency.NewSession(client)
	if err != nil {
		return nil, err
	}
	lock.m = concurrency.NewMutex(lock.s, key)
	return lock, nil
}
