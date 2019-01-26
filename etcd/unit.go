/********************************************
author:
  nike：freetoo
  name: yigui-lu
 wx/qq: 48092788
e-mail: gcode@qq.com
*********************************************/

package etcd

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"strings"

	//"time"

	etcd3 "github.com/coreos/etcd/clientv3"
)

func NewConf(target, caFile, certFile, keyFile string) (*etcd3.Config, error) {
	var tlsCfg *tls.Config = nil
	if (caFile != "") && (certFile != "") && (keyFile != "") {
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
	}
	return &etcd3.Config{
		Endpoints: strings.Split(target, ","),
		//DialTimeout: 5 * time.Second,
		TLS: tlsCfg,
	}, nil
}
