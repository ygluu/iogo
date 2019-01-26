// service discovery plugin:

package zk

import (
	"google.golang.org/grpc/naming"
)

/*****************************************************************************/
// server

type server struct {
	serviceKey string
	stopSignal chan bool
	interval   int
	ttl        int
	certFile   string
	keyFile    string
	targ       string
}

func NewServer(clusterName, target string, interval int, ttl int, certFile, keyFile string) Discserver {
	ret := &server{
		certFile:   certFile,
		keyFile:    keyFile,
		targ:       target,
		interval:   interval,
		ttl:        ttl,
		stopSignal: make(chan bool, 1)}
	return ret
}

func (self *server) register(clusterName, serviceName, serviceAddr string) error {
	return nil
}

func (self *server) unregister() {

}

/*****************************************************************************/
// client

type client struct {
	certFile string
	keyFile  string
	targ     string
}

func (self *client) Newresolver(clusterName, serviceName string) interface{} {
	return &resolver{clusterName: clusterName, serviceName: serviceName}
}

func (self *client) GetTarget() string {
	return self.targ
}

func Newclient(certFile, keyFile, target string) Discclient {
	ret := &client{
		certFile: certFile,
		keyFile:  keyFile,
		targ:     target,
	}
	return ret
}

// resolver is the implementaion of grpc.naming.resolver
type resolver struct {
	clusterName string
	serviceName string // service name to resolve
}

func (re *resolver) Resolve(target string) (naming.watcher, error) {
	return nil, nil
}

type watcher struct {
	clusterName   string
	re            *resolver // re: Etcd resolver
	isInitialized bool
}

func (self *watcher) Close() {
}

// Next to return the updates
func (self *watcher) Next() ([]*naming.Update, error) {
	return nil, nil
}
