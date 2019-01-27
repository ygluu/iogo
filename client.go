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
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/naming"
)

/*****************************************************************************/
// grpc client
type grpcConnect struct {
	owner *grpcClient
	cli   *grpc.ClientConn
}

type grpcClient struct {
	clusterName string
	conns       []*grpcConnect
	disc        DiscClient
}

func (self *grpcClient) start() {
	log.Printf("iogo: client: waitfor inline of service discovery(%s)......\n", self.disc.GetTarget())
	for {
		if self.disc.IsInline() {
			break
		}
		time.Sleep(1000)
	}
	time.Sleep(time.Second * 2)
	log.Printf("iogo: client: service discovery inline")
}

func (self *grpcClient) close() {
	for _, conn := range self.conns {
		conn.cli.Close()
	}
}

func (self *grpcClient) newConnect(ctx context.Context, serviceName string, opts ...grpc.DialOption) *grpc.ClientConn {
	conn := &grpcConnect{
		owner: self,
	}
	self.conns = append(self.conns, conn)
	r := self.disc.NewResolver(self.clusterName, serviceName)
	b := grpc.RoundRobin(r.(naming.Resolver))
	if ctx == nil {
		ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	}
	var err error
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBalancer(b))
	conn.cli, err = grpc.DialContext(ctx, self.disc.GetTarget(), opts...)
	if err != nil {
		log.Printf("error: iogo.newConnect.DialContext: %v\n", err)
		return nil
	}
	return conn.cli
}

func (self *grpcClient) newConnectByAddr(ctx context.Context, addr string, opts ...grpc.DialOption) *grpc.ClientConn {
	conn := &grpcConnect{
		owner: self,
	}
	if ctx == nil {
		ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	}
	var err error
	opts = append(opts, grpc.WithInsecure())
	conn.cli, err = grpc.DialContext(ctx, addr, opts...)
	if err != nil {
		log.Printf("error: iogo.newConnectByAddr.DialContext: %v\n", err)
		return nil
	}
	return conn.cli
}

var client *grpcClient = nil

func NewClient(clusterName string, disc DiscClient) *grpcClient {
	if client == nil {
		client = &grpcClient{
			clusterName: clusterName,
			disc:        disc,
		}
	}
	return client
}

func NewClientCreds(serverName string, caFile, certFile, keyFile string) grpc.DialOption {
	if (caFile != "") && (certFile != "") && (keyFile != "") {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Fatalf("error: iogo.NewClientCreds.LoadX509KeyPair: %v", err)
		}
		var creds credentials.TransportCredentials
		if caFile != "" {
			certPool := x509.NewCertPool()
			ca, err := ioutil.ReadFile(caFile)
			if err != nil {
				log.Fatalf("error: iogo.NewClientCreds.ReadFile err: %v", err)
			}

			if ok := certPool.AppendCertsFromPEM(ca); !ok {
				log.Fatalf("error: iogo.NewClientCreds.certPool.AppendCertsFromPEM err")
			}

			creds = credentials.NewTLS(&tls.Config{
				Certificates: []tls.Certificate{cert},
				ServerName:   serverName,
				RootCAs:      certPool,
			})
		} else {
			creds = credentials.NewTLS(&tls.Config{
				Certificates: []tls.Certificate{cert},
				ServerName:   serverName,
			})
		}

		return grpc.WithTransportCredentials(creds)
	} else if caFile != "" {
		creds, err := credentials.NewClientTLSFromFile(caFile, serverName)
		if err != nil {
			log.Fatalf("error: iogo.NewClientCreds.NewClientTLSFromFile: %v", err)
		}
		return grpc.WithTransportCredentials(creds)
	} else {
		//log.Fatalf("error: iogo.NewClientCreds: param values")
		return nil
	}
}

func NewConnect(ctx context.Context, serviceName string, opts ...grpc.DialOption) *grpc.ClientConn {
	if client == nil {
		log.Fatalf("error: client no init")
	}
	return client.newConnect(ctx, serviceName, opts...)
}

func NewConnectByAddr(ctx context.Context, Addr string, opts ...grpc.DialOption) *grpc.ClientConn {
	// Connection pooling should be used here
	if client == nil {

		log.Fatalf("error: client no init")
	}
	return client.newConnectByAddr(ctx, Addr, opts...)
}

func startCli() {
	if client != nil {
		log.Printf("iogo: client start......")
		client.start()
	}
}

func stopCli() {
	if client != nil {
		client.close()
		log.Printf("iogo: client close")
	}
}
