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

package hash // import "iogo/plugin/lb/hash"

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
	"time"

	"iogo/plugin"
)

type hashRing []uint32

func (s hashRing) Len() int {
	return len(s)
}

func (s hashRing) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s hashRing) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type addrInfo struct {
	weight    int
	addr      string
	connected bool
	metadata  interface{}
}

type addrNode struct {
	next      int
	addrInfos []*addrInfo
}

func (self *addrNode) get() string {
	if len(self.addrInfos) == 0 {
		return ""
	}
	count := len(self.addrInfos)
	for count > 0 {
		count--
		ai := self.addrInfos[self.next]
		self.next++
		if self.next >= len(self.addrInfos) {
			self.next = 0
		}
		if ai.connected {
			return ai.addr
		}
	}
	return ""
}

func (self *addrNode) add(ai *addrInfo) {
	for _, v := range self.addrInfos {
		if v.addr == ai.addr {
			return
		}
	}
	self.addrInfos = append(self.addrInfos, ai)
}

func (self *addrNode) del(addr string) bool {
	for i, v := range self.addrInfos {
		if v.addr == addr {
			self.addrInfos = append(self.addrInfos[:i], self.addrInfos[i+1:]...)
			return len(self.addrInfos) == 0
		}
	}
	return false
}

func (self *addrNode) exsit(addr string) *addrInfo {
	for _, v := range self.addrInfos {
		if v.addr == addr {
			return v
		}
	}
	return nil
}

type HashBalan struct {
	replicas int
	hashRing hashRing
	hashMap  map[uint32]*addrNode
	addrMap  map[string]*addrInfo
	getCount uint64
	Mutex    sync.Mutex
	addrCh   chan int
}

func (self *HashBalan) Clone() plugin.LoadBalan {
	return NewLoadBalan(self.replicas)
}

func (self *HashBalan) IsEmpty() bool {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	return len(self.addrMap) == 0
}

func (self *HashBalan) Exsit(addr string) bool {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	flag := self.addrMap[addr]
	return flag != nil
}

func (self *HashBalan) sortHashRing() {
	self.hashRing = hashRing{}
	for k := range self.hashMap {
		self.hashRing = append(self.hashRing, k)
	}
	sort.Sort(self.hashRing)
}

func (self *HashBalan) getHashByCount(i uint64, addr string) uint32 {
	str := strconv.FormatUint(i, 10) + addr
	data := md5.Sum([]byte(str))
	str = fmt.Sprintf("%x", data)
	return crc32.ChecksumIEEE([]byte(str))
}

func (self *HashBalan) getHash(addr string) uint32 {
	data := md5.Sum([]byte(addr))
	addr = fmt.Sprintf("%x", data)
	return crc32.ChecksumIEEE([]byte(addr))
}

func (self *HashBalan) Add(addr string, weight int) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	ai := self.addrMap[addr]
	if ai != nil {
		if ai.weight == weight {
			return
		}
		self.Del(addr)
	}

	if weight < 1 {
		weight = 1
	}
	w := weight * self.replicas
	ai = &addrInfo{
		addr:      addr,
		weight:    weight,
		connected: true,
	}
	self.addrMap[addr] = ai

	for i := 0; i < w; i++ {
		hash := self.getHashByCount(uint64(i), addr)
		node := self.hashMap[hash]
		if node == nil {
			node = &addrNode{next: 0}
			self.hashMap[hash] = node
		}
		node.add(ai)
	}
	self.sortHashRing()

	select {
	case <-self.addrCh:
	default:
	}
	self.addrCh <- 1
}

func (self *HashBalan) Del(addr string) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	ai := self.addrMap[addr]
	if ai == nil {
		return
	}
	delete(self.addrMap, addr)
	w := ai.weight * self.replicas
	for i := 0; i < w; i++ {
		hash := self.getHashByCount(uint64(i), addr)
		node := self.hashMap[hash]
		if node == nil {
			continue
		}
		if !node.del(addr) {
			continue
		}
		delete(self.hashMap, hash)
	}
	self.sortHashRing()
	select {
	case <-self.addrCh:
	default:
	}
	self.addrCh <- 1
}

func (self *HashBalan) Get(key string) string {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	if len(self.addrMap) == 0 {
		return ""
	}
	for i := 0; i < 3; i++ {
		var hash uint32
		if key == "" {
			self.getCount++
			hash = self.getHashByCount(self.getCount, strconv.FormatUint(uint64(time.Now().UnixNano()), 10))
		} else {
			hash = self.getHash(key)
		}
		index := sort.Search(len(self.hashRing), func(i int) bool { return self.hashRing[i] >= hash })
		if index >= len(self.hashRing) {
			index = len(self.hashRing) - 1
		} else if index == len(self.hashRing)-1 {
			index = 0
		}
		hash = self.hashRing[index]
		ai := self.hashMap[hash]
		ret := ai.get()
		if ret != "" {
			return ret
		}
	}
	return ""
}

func (self *HashBalan) Up(addr string, metadata interface{}) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	ai := self.addrMap[addr]
	if ai != nil {
		ai.connected = true
		ai.metadata = metadata
	}
}

func (self *HashBalan) Down(addr string) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	ai := self.addrMap[addr]
	if ai != nil {
		ai.connected = false
	}
}

func (self *HashBalan) Notify() <-chan int {
	return self.addrCh
}

func (self *HashBalan) Address() []plugin.Address {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	ret := make([]plugin.Address, len(self.addrMap))
	i := 0
	for _, v := range self.addrMap {
		ret[i].Addr = v.addr
		ret[i].Metadata = v.metadata
		i++
	}
	return ret
}

func NewLoadBalan(replicas int) plugin.LoadBalan {
	if replicas <= 0 {
		replicas = 150
	}
	ret := &HashBalan{
		replicas: replicas,
		hashMap:  make(map[uint32]*addrNode),
		addrMap:  make(map[string]*addrInfo),
		addrCh:   make(chan int, 1),
	}
	return ret
}
