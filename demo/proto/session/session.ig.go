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

var CallSession SessionClient = nil
var CallAccount AccountClient = nil

func onStart() {

}

func init() {
	iogo.RegisterFuncOnStart(onStart)
	
	cliSession := &sessionClient{}
	CallSession = cliSession
	iogo.RegisterClient("session.Session.proto", "proto.Session", CallSession, "")
	
	cliAccount := &accountClient{}
	CallAccount = cliAccount
	iogo.RegisterClient("session.Session.proto", "proto.Account", CallAccount, "")
}

func RegisterSessionServer(srv SessionServer) {
	iogo.RegisterServer("session.Session.proto", "proto.Session", srv, &_Session_serviceDesc)
}

func RegisterAccountServer(srv AccountServer) {
	iogo.RegisterServer("session.Session.proto", "proto.Account", srv, &_Account_serviceDesc)
}
