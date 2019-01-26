package main

var confFile string = `package config

import (
	"flag"
	"fmt"
	"net"
	"os"
	"path"
	"strings"
)

const (
	maskIpDef = "192.168.1.0"
	etcdsDef = ""
	//etcdsDef = "http://192.168.137.128:2379"
)

var (
	ClusterName    string = "iogo-demo"
	SvcName		   string = ""
	SvcNameHello   string = "svc-Hello"
	SvcAddr        string = ""
	SvcCaFile      string = ""
	SvcCertFile    string = ""
	SvcKeyFile     string = ""
	CliCaFile      string = ""
	CliCertFile    string = ""
	CliKeyFile     string = ""
	Etcds          string = ""
	EtcdCaFile     string = ""
	EtcdCertFile   string = ""
	EtcdKeyFile    string = ""
)

func getExeName() string {
	ss := strings.Split(path.Base(strings.Replace(os.Args[0], "\\", "/", -1)), ".")
	return ss[0]
}

func checkIp(ip, maskIp string) bool {
	ip1 := strings.Split(ip, ".")
	ip2 := strings.Split(maskIp, ".")
	ret := true
	for i, item := range ip2 {
		if item == "0" {
			break
		} else if ip1[i] != ip2[i] {
			ret = false
			break
		}
	}
	return ret
}

func getAddr(maskIp string) string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if (ipnet.IP.To4() != nil) && (checkIp(ipnet.IP.String(), maskIp)) {
				return ipnet.IP.String()
			}

		}
	}
	return ""
}

func checkAddr(ip string, port int) string {
	ret := fmt.Sprintf("%s:%d", ip, port)
	for {
		lis, err := net.Listen("tcp", ret)
		if err == nil {
			lis.Close()
			return ret
		}
		port++
		ret = fmt.Sprintf("%s:%d", ip, port)
	}
	return ret
}

func getEtcdAddr(ip string) string {
	return fmt.Sprintf("http://%s:2379", ip)
}

// cmd line input or config file load...
func init() {
	var (
		clus         = flag.String("cluster", ClusterName, "cluster name")
		serv         = flag.String("svcname", getExeName(), "service name")
		mask         = flag.String("maskip", maskIpDef, "mask ip of listen")
		port         = flag.Int("svcport", 20001, "listening port")
		etcds        = flag.String("etcds", etcdsDef, "etcd cluster address: http://ip1:port1,http://ip2:port2, default: LocalIP:2379")
		svcCaFile    = flag.String("svccafile", "", "ca file name of service")
		svcCertFile  = flag.String("svccertfile", "", "cert file name of service")
		svcKeyFile   = flag.String("svckeyfile", "", "key file name of service")
		etcdCaFile   = flag.String("etcdcafile", "", "ca file name of etcd cluster")
		etcdCertFile = flag.String("etcdcertfile", "", "cert file name of etcd cluster")
		etcdKeyFile  = flag.String("etcdkeyfile", "", "key file name of etcd cluster")
		cliCaFile    = flag.String("clicafile", "", "ca file name of client")
		cliCertFile  = flag.String("clicertfile", "", "cert file name of client")
		cliKeyFile   = flag.String("clikeyfile", "", "key file name of client")
	)
	flag.Parse()

	ClusterName = *clus
	SvcName = *serv
	localIP := getAddr(*mask)
	SvcAddr = checkAddr(localIP, *port)
	if *etcds == "" {
		Etcds = getEtcdAddr(localIP)
	} else {
		Etcds = *etcds
	}
	SvcCaFile = *svcCaFile
	SvcCertFile = *svcCertFile
	SvcKeyFile = *svcKeyFile
	CliCaFile = *cliCaFile
	CliCertFile = *cliCertFile
	CliKeyFile = *cliKeyFile
	EtcdCaFile = *etcdCaFile
	EtcdCertFile = *etcdCertFile
	EtcdKeyFile = *etcdKeyFile
}

`

var iogoFile string = `package iogo

import (
	"fmt"
	"log"

	etcd3 "github.com/coreos/etcd/clientv3"
	"google.golang.org/grpc"

	offi_iogo "iogo"
	"iogo/etcd"

	"iogo-demo/config"
)

var etcdConf *etcd3.Config = nil

func init() {
	etcdConf, err := etcd.NewConf(config.Etcds, config.EtcdCaFile, config.EtcdCertFile, config.EtcdKeyFile)
	if err != nil {
		log.Fatalf("iogo-app: GetEtcdCfg:%v", err)
	}
	cliEtcd := etcd.NewClient(etcdConf)
	offi_iogo.NewClient(config.ClusterName, cliEtcd)

	svrEtcd := etcd.NewServer(10, 15, etcdConf)
	creds := offi_iogo.NewServerCreds(config.SvcCaFile, config.SvcCertFile, config.SvcKeyFile)
	if creds != nil {
		offi_iogo.NewServer(config.ClusterName, config.SvcName, config.SvcAddr, svrEtcd, creds)
	} else {
		offi_iogo.NewServer(config.ClusterName, config.SvcName, config.SvcAddr, svrEtcd)
	}
}

func GetServerHandle() *grpc.Server {
	return offi_iogo.GetServerHandle()
}

func NewConnect(serviceName string) *grpc.ClientConn {
	creds := offi_iogo.NewClientCreds(config.SvcName, config.CliCaFile, config.CliCertFile, config.CliKeyFile)
	if creds != nil {
		return offi_iogo.NewConnect(nil, serviceName, creds)
	} else {
		return offi_iogo.NewConnect(nil, serviceName)
	}
}

func NewLock(key string) (*etcd.Lock, error) {
	return etcd.NewLock(fmt.Sprintf("/%s/locks/%s", config.ClusterName, key), etcdConf)
}

func Start() {
	log.Printf("The service:%s of cluster:%s  is starting...", config.SvcName, config.ClusterName)
	offi_iogo.Start()
}

func Stop() {
	offi_iogo.Stop()
	log.Printf("The service:%s of cluster:%s is stop", config.SvcName, config.ClusterName)
}

`

var protoFile string = `syntax = "proto3";

option objc_class_prefix = "HLW";

package proto;

service Hello {
	rpc SayHello(SayHelloRequest) returns (SayHelloReply) {	
	}
}

message SayHelloRequest {
    string Text = 1;
	int32 Flag = 2;
}

message SayHelloReply {
    string Text = 1;
	int32 Flag = 2;
}

`

var mainFile string = `package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"iogo-demo/iogo"
	_ "iogo-demo/svc-demo/service"
)

func waitStop() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		log.Printf("receive signal '%v'", s)
		iogo.Stop()
		os.Exit(1)
	}()
}

func main() {
	waitStop()

	iogo.Start()
}

`

var clientFile string = `package service

import (
	"iogo-demo/config"
	"iogo-demo/iogo"
	"iogo-demo/proto"
)

type cli struct {
	hello proto.HelloClient
	// ...
}

var client *cli = nil

func init() {
	client = &cli{}
	
	// demo connect Hello
	c := iogo.NewConnect(config.SvcNameHello)
	client.hello = proto.NewHelloClient(c)
	// ...

}

`
var serverFile string = `package service

import (
	"log"

	"golang.org/x/net/context"

	"iogo"
	"iogo-demo/proto"
)

type server struct {
}

func init() {
	proto.RegisterHelloServer(iogo.GetServerHandle(), &server{})
}

func (self *server) SayHello(ctx context.Context, request *proto.SayHelloRequest) (reply *proto.SayHelloReply, err error) {
	log.Printf("client say: %s (flag:%d)\n", request.Text, request.Flag)

	reply = &proto.SayHelloReply{
		Text: "Welcome to iogo",
		Flag: request.Flag,
	}

	return reply, nil
}

`

var testCliFile string = `package main

import (
	"log"

	"iogo"
	"iogo/etcd"

	"iogo-demo/config"
	"iogo-demo/proto"
)

type cli struct {
	hello proto.HelloClient
	// ...
}

var client *cli = nil

func init() {
	client = &cli{}

	etcdConf, err := etcd.NewConf(config.Etcds, "", "", "")
	if err != nil {
		log.Fatalf("iogo-app: NewConf:%v", err)
	}
	cliEtcd := etcd.NewClient(etcdConf)
	iogo.NewClient(config.ClusterName, cliEtcd)
	iogo.Start()

	c := iogo.NewConnect(nil, config.SvcNameHello)
	client.hello = proto.NewHelloClient(c)
	// ...

}

`
var testMainFile string = `package main

import (
	"log"
	"time"

	"golang.org/x/net/context"

	"iogo-demo/proto"
)

func main() {

	ticker := time.NewTicker(1 * time.Second)
	for t := range ticker.C {
		req := &proto.SayHelloRequest{
			Text: "Hello World",
			Flag: int32(t.Second()),
		}
		log.Printf("Request: %s (falg:%d)\n", req.Text, req.Flag)
		res, err := client.hello.SayHello(context.Background(), req)
		if err == nil {
			log.Printf("Reply: %s (falg:%d)\n", res.Text, res.Flag)
		} else {
			log.Println(err)
		}
	}
}

`
