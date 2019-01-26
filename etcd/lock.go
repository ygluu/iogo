/********************************************
author:
  nike：freetoo
  name: yigui-lu
 wx/qq: 48092788
e-mail: gcode@qq.com
*********************************************/

package etcd

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
