package iogo

import (
	"iogo"
	"iogo/plugin"
	"iogo/plugin/cdc/json"
	"iogo/plugin/cdc/pb"
	"iogo/plugin/lb/hash"
	"iogo/plugin/net/http"
	"iogo/plugin/sd/etcd"
)

// For service discovery
const defClusterName = "IogoCluster"

// endpoints of Service discovery
const defSdEndpoints = "http://192.168.1.23:2379"

// Select IP based on subnet mask
const defSubnetMask = "192.168.1.2/22"

func getNewNetOption(addr, flag string) (svr http.NewOption, cli http.NewOption) {
	if flag == "mutway" {
		caFile := iogo.CfgReader.GetCaFileName("../keys/mutway/ca.crt")
		svrCertFile := iogo.CfgReader.GetServerCertFileName("../keys/mutway/server.crt")
		svrKeyFile := iogo.CfgReader.GetServerKeyFileName("../keys/mutway/server.key")
		cliCertFile := iogo.CfgReader.GetClientCertFileName("../keys/mutway/client.crt")
		cliKeyFile := iogo.CfgReader.GetClientKeyFileName("../keys/mutway/client.key")
		certSvrName := iogo.CfgReader.GetCertServerName("")
		cli := http.WithCertFilesCli(cliCertFile, cliKeyFile, caFile, certSvrName)
		svr := http.WithCertFilesSvr(addr, svrCertFile, svrKeyFile, caFile)
		return svr, cli
	} else if flag == "oneway" {
		caFile := ""
		svrCertFile := iogo.CfgReader.GetServerCertFileName("../keys/oneway/server.crt")
		svrKeyFile := iogo.CfgReader.GetServerKeyFileName("../keys/oneway/server.key")
		cliCertFile := svrCertFile
		cliKeyFile := ""
		certSvrName := ""
		cli := http.WithCertFilesCli(cliCertFile, cliKeyFile, caFile, certSvrName)
		svr := http.WithCertFilesSvr(addr, svrCertFile, svrKeyFile, caFile)
		return svr, cli
	} else {
		return http.WithAddr(addr), http.WithDefault()
	}
}

func getNewEtcdOption(endpoints, flag string) etcd.NewOption {
	if flag == "mutway" {
		caFile := iogo.CfgReader.GetSdCaFileName("../keys/mutway/ca.crt")
		certFile := iogo.CfgReader.GetSdCertFileName("../keys/mutway/client.crt")
		keyFile := iogo.CfgReader.GetSdKeyFileName("../keys/mutway/client.key")
		certSvrName := iogo.CfgReader.GetCertServerName("")
		return etcd.WithCertFile(endpoints, certFile, keyFile, caFile, certSvrName)
	} else if flag == "oneway" {
		caFile := ""
		certFile := iogo.CfgReader.GetSdCertFileName("../keys/oneway/server.crt")
		keyFile := iogo.CfgReader.GetSdKeyFileName("../keys/oneway/server.key")
		certSvrName := iogo.CfgReader.GetCertServerName("")
		return etcd.WithCertFile(endpoints, certFile, keyFile, caFile, certSvrName)
	} else {
		return etcd.WithEndpoints(endpoints)
	}
}

func getCodec(flag string) plugin.Codec {
	if flag == "json" {
		return json.NewCodec()
	} else {
		return pb.NewCodec()
	}
}

func RunIogo(servicePort string) {
	sdEndpoints := iogo.CfgReader.GetSdEndpoints(defSdEndpoints)
	clusterName := iogo.CfgReader.GetClusterName(defClusterName)
	serviceAddr := iogo.CfgReader.GetServiceAddr("", servicePort, defSubnetMask)

	newOptEtcd := getNewEtcdOption(sdEndpoints, "")                // or mutway stl, or oneway  stl, or "" not saft
	newOptNetSvr, newOptNetCli := getNewNetOption(serviceAddr, "") // or mutway stl, or oneway  stl, or "" not saft

	pluginCli := iogo.NewPluginClient()
	pluginCli.Cdc = getCodec("json") // or pb, or json
	pluginCli.Lb = hash.NewLoadBalan(200)
	pluginCli.Sd = etcd.NewClient(newOptEtcd)
	pluginCli.Net = http.NewClient(newOptNetCli)

	pluginSvr := iogo.NewPluginServer()
	if servicePort != "" {
		pluginSvr.Cdc = pluginCli.Cdc
		pluginSvr.Net = http.NewServer(newOptNetSvr)
		pluginSvr.Sd = etcd.NewServer(newOptEtcd, 20, 25)
	}

	iogo.SetPlugin(pluginSvr, pluginCli)
	iogo.Run(clusterName, serviceAddr)
}

func StopIogo() {
	iogo.Stop()
}
