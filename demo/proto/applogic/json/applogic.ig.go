package proto

import "iogo"

import (
	"context"
)

type Context = context.Context

func TODO() context.Context {
	return context.TODO()
}

func Background() context.Context {
	return context.Background()
}

var CallMethod iogo.FuncCallService = nil

type applogicClient struct {
}

func (self *applogicClient) Say(ctx Context, request *SayRequest) (*SayReply, error) {
	reply := &SayReply{}
	return reply, CallMethod(ctx, request, reply, "")
}

var CallApplogic *applogicClient = &applogicClient{}

func onStart() {

}

func init() {
	iogo.RegisterFuncOnStart(onStart)
	
	iogo.RegisterClient("json.applogic.go.proto", "proto.Applogic", "", &CallMethod)
}

type ApplogicServer interface {
	Say(ctx Context, request *SayRequest) (*SayReply, error)
}

func RegisterApplogicServer(srv ApplogicServer) {
	iogo.RegisterServer("json.applogic.go.proto", "proto.Applogic", srv, nil)
}
