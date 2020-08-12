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

package tls // import "iogo/unit/tls"

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
)

func NewTlsServer(certFile, keyFile, caFile string) (*tls.Config, error) {
	if certFile == "" || keyFile == "" {
		return nil, fmt.Errorf("NewTlsServer: Certfile or keyfile file name cannot be empty")
	}

	cfg := &tls.Config{}

	if caFile != "" {
		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(caFile)
		if err != nil {
			return nil, fmt.Errorf("NewTlsServer: Read file error: %v file:%s", err, caFile)
		}
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			return nil, fmt.Errorf("NewTlsServer: AppendCertsFromPEM error file:%s", caFile)
		}
		cfg.ClientCAs = certPool
		cfg.ClientAuth = tls.RequireAndVerifyClientCert
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("NewTlsServer: tls.LoadX509KeyPair err: %v certfile:%s keyfile:%s", err, certFile, keyFile)
	}
	cfg.Certificates = []tls.Certificate{cert}
	return cfg, nil
}

func NewTlsClient(certFile, keyFile, caFile, certSvrName string) (*tls.Config, error) {
	if certFile == "" {
		return nil, fmt.Errorf("NewTlsClient: Certfile file name cannot be empty")
	}

	cfg := &tls.Config{
		ServerName:         certSvrName,
		InsecureSkipVerify: certSvrName == "",
	}

	if keyFile != "" && caFile != "" {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, fmt.Errorf("NewTlsClient: tls.LoadX509KeyPair err: %v certfile:%s keyfile:%s", err, certFile, keyFile)
		}
		cfg.Certificates = []tls.Certificate{cert}
	}

	if caFile == "" {
		caFile = certFile
	}
	ca, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("NewTlsClient: Read file error: %v file:%s", err, certFile)
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(ca) {
		return nil, fmt.Errorf("NewTlsClient: AppendCertsFromPEM error file:%s", caFile)
	}
	cfg.ClientCAs = certPool

	return cfg, nil
}
