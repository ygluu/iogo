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

package http // import "iogo/plugin/net/http"

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
	"unsafe"

	"iogo/log"
	"iogo/plugin"
	"iogo/unit"
	"iogo/unit/md"
	"iogo/unit/tls"
)

type svcClientInfo struct {
	owner       *client
	serviceName string
	target      string
	matadata    string
	lb          plugin.LoadBalan
	codec       plugin.Codec
}

func (self *svcClientInfo) CallService(fullMethod string, ctx context.Context, request interface{}, response interface{}) error {
	codec := self.owner.codec
	rJson, err := codec.Marshal(request)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", "", strings.NewReader(string(rJson)))
	if err != nil {
		return err
	}
	md, _ := md.FromOutgoingContext(ctx)
	for k, v := range md {
		req.Header[k] = v
	}

	getUrl := func(svcAddr string) string {
		var http string
		if self.owner.tr.TLSClientConfig == nil {
			http = "http"
		} else {
			http = "https"
		}
		ret := fmt.Sprintf("%s://%s%s", http, svcAddr, fullMethod)
		return ret
	}

	client := self.owner.newClient()
	client.Timeout = 15 * time.Second

	var resp *http.Response
	for i := 0; i < 3; i++ {
		var svcAddr string
		svcAddr = self.lb.Get("")
		if svcAddr == "" {
			return fmt.Errorf("Get service addr fail")
		}
		url, _ := url.Parse(getUrl(svcAddr))
		req.URL = url
		resp, err = client.Do(req)
		if err == nil {
			break
		}
	}
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	eStr := resp.Header.Get(IOGO_ERROR_NAME)
	respBody, err := ioutil.ReadAll(resp.Body)
	if len(respBody) == 0 {
		return fmt.Errorf("Call service not response")
	}
	if eStr == IOGO_ERROR_VALUE {
		return fmt.Errorf(string(respBody))
	}
	err = codec.Unmarshal(respBody, response)
	if err != nil {
		return fmt.Errorf(string(respBody))
	}

	return nil
}

type client struct {
	tr      *HttpTransport
	cli     *HttpClient
	lb      plugin.LoadBalan
	codec   plugin.Codec
	isFirst bool
}

func (self *client) newClient() *http.Client {
	tr := &HttpTransport{}
	*tr = *self.tr
	ret := &http.Client{}
	*ret = *self.cli
	ret.Transport = tr
	return ret
}

func (self *client) SetCodec(codec plugin.Codec) {
	self.codec = codec
}

func (self *client) Register(matadata, serviceName string, svcClient interface{}, loadBalan plugin.LoadBalan) {
	if !unit.IsStructPoint(svcClient) {
		panic("RegisterClient: Parameter svcClient must be a structure pointer")
	}
	var svcPFunc *plugin.FuncCallHttpSvc = nil
	v := reflect.ValueOf(svcClient).Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := t.Field(i)
		if f.Type.String() == "plugin.FuncCallHttpSvc" {
			svcPFunc = (*plugin.FuncCallHttpSvc)(unsafe.Pointer(v.Field(i).Addr().Pointer()))
		}
	}
	if svcPFunc == nil {
		panic(fmt.Sprintf("HTTP plug-ins cannot use other source code: %s", t.String()))
	}

	self.lb = loadBalan
	svcc := &svcClientInfo{
		owner:       self,
		serviceName: serviceName,
		lb:          loadBalan,
		matadata:    matadata,
		codec:       self.codec,
	}
	*svcPFunc = svcc.CallService
	if self.isFirst {
		self.isFirst = false
		log.I(LC_HTTPC_READY, "HttpClient: Client is ready")
	}
}

func NewClient(newOpt NewOption) plugin.NetClient {
	var cli *HttpClient
	var tr *HttpTransport
	if newOpt.Flag() == flag_client {
		cli = newOpt.Matadata().(*HttpClient)
		v := reflect.ValueOf(cli.Transport).Elem().Addr()
		tr = v.Interface().(*HttpTransport)
	} else if newOpt.Flag() == flag_tlscfg {
		cfg := newOpt.Matadata().(*CfgTls)
		cfg.InsecureSkipVerify = true
		tr := &http.Transport{TLSClientConfig: cfg}
		cli = &http.Client{Transport: tr}
	} else if newOpt.Flag() == flag_certfilecli {
		md := newOpt.Matadata().(map[string]string)
		certSvrName := md["certsvrname"]
		caFile := md["cafile"]
		certFile := md["certfile"]
		keyFile := md["keyfile"]
		cfg, err := tls.NewTlsClient(certFile, keyFile, caFile, certSvrName)
		if err != nil {
			panic("NewHttpClient: cert file error:" + err.Error())
		}
		tr = &http.Transport{TLSClientConfig: cfg}
		cli = &http.Client{Transport: tr}
	} else if newOpt.Flag() == flag_trans {
		tr = newOpt.Matadata().(*HttpTransport)
		cli = &http.Client{Transport: tr}
	} else if newOpt.Flag() == flag_default {
		tr = &http.Transport{}
		cli = &http.Client{Transport: tr}
	} else {
		//t := reflect.TypeOf(newOpt.Matadata())
		panic("NewHttpClient: NewOption error: " + newOpt.Flag())
	}
	return &client{tr: tr, cli: cli, isFirst: true}
}
