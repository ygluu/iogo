package proto

import "context"

import "iogo"
import "iogo/unit/md"

type Context = context.Context

func TODO() context.Context {
	return context.TODO()
}

func Background() context.Context {
	return context.Background()
}

func TokenContext(token string) context.Context {
	ctx := TODO()
	ctx = md.AppendToOutgoingContext(ctx, "set-token", token)
	return ctx
}

var CallApplogic ApplogicClient = nil

func onStart() {

}

func init() {
	iogo.RegisterFuncOnStart(onStart)
	
	cliApplogic := &applogicClient{}
	CallApplogic = cliApplogic
	iogo.RegisterClient("applogic.applogic.proto", "proto.Applogic", CallApplogic, "")
}

func RegisterApplogicServer(srv ApplogicServer) {
	iogo.RegisterServer("applogic.applogic.proto", "proto.Applogic", srv, &_Applogic_serviceDesc)
}
