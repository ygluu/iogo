package etcd

import (
	"bytes"
	"iogo/etcd"
	"log"
	"runtime"
	"strconv"
	"testing"
	"time"
)

const key = "/etcdlock-test1/"

func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func do_test(c chan bool) {
	conf, err := etcd.NewConf("http://192.168.137.128:2379", "", "", "")
	if err != nil {
		log.Println(GetGID(), err)
		c <- true
		return
	}
	lock, err := etcd.NewLock(key, conf)
	if err != nil {
		log.Println(GetGID(), err)
		c <- true
		return
	}
	for i := 0; i < 100; i++ {
		log.Println(GetGID(), "do lock")
		err = lock.Lock()
		if err != nil {
			log.Println(GetGID(), err)
			continue
		}
		log.Println(GetGID(), "lock")
		time.Sleep(time.Second)
		lock.Unlock()
		log.Println(GetGID(), "unlock")
	}
	c <- true
}

func Start(t *testing.T) {
	log.Println("start test")

	c := make(chan bool, 2)
	go do_test(c)
	go do_test(c)
	<-c
	<-c

	log.Println("end test")
}
