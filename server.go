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

// communication protocole plugin: http2

package iogo

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

/*****************************************************************************/
// grpc Server
type grpcServer struct {
	clusterName string
	serviceName string
	svcAddr     string
	disc        DiscServer
	svr         *grpc.Server
}

func (self *grpcServer) start() {
	lis, err := net.Listen("tcp", self.svcAddr)
	if err != nil {
		panic(err)
	}
	regfunc := func() {
		time.Sleep(time.Second * 2)
		self.disc.Register(self.clusterName, self.serviceName, self.svcAddr)
	}
	go regfunc()
	log.Printf("iogo: service %s(%s) start ok", self.serviceName, self.svcAddr)
	self.svr.Serve(lis)
}

func (self *grpcServer) stop() {
	self.disc.Unregister()
	log.Printf("iogo: service %s(%s) stop ok", self.serviceName, self.svcAddr)
	self.svr.Stop()
}

func (self *grpcServer) getHandle() *grpc.Server {
	return self.svr
}

func NewServerCreds(caFile, certFile, keyFile string) grpc.ServerOption {
	if (caFile != "") && (certFile != "") && (keyFile != "") {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Fatalf("error: iogo.NewServerCredentials.LoadX509KeyPair : %v", err)
		}

		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(caFile)
		if err != nil {
			log.Fatalf("error: iogo.NewServerCredentials.ReadFile err: %v", err)
		}

		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			log.Fatalf("error: iogo.NewCredentials.certPool.AppendCertsFromPEM err")
		}

		creds := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    certPool,
		})
		return grpc.Creds(creds)
	} else if (certFile != "") && (keyFile != "") {
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			log.Fatalf("error: iogo.NewCredentials.NewServerTLSFromFile: %v", err)
		}
		return grpc.Creds(creds)
	} else {
		//log.Fatalf("error: iogo.NewCredentials: param values")
		return nil
	}
}

var server *grpcServer = nil

func GetServerHandle() *grpc.Server {
	if server == nil {
		log.Fatalf("error: server no init")
	}
	return server.getHandle()
}

func NewServer(clusterName string, serviceName string, svcAddr string, disc DiscServer, opt ...grpc.ServerOption) *grpcServer {
	if server == nil {
		server = &grpcServer{
			clusterName: clusterName,
			disc:        disc,
			serviceName: serviceName,
			svcAddr:     svcAddr,
			svr:         grpc.NewServer(opt...),
		}
	}
	return server
}

func startSvr() {
	if server != nil {
		log.Printf("iogo: server start...")
		server.start()
	}
}

func stopSvr() {
	if server != nil {
		server.stop()
		log.Printf("iogo: server stop")
	}
}
