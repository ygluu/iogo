package service

import (
	"iogo/demo/comm/token"
	"iogo/demo/proto/session"
)

type server struct {
}

func init() {
	proto.RegisterSessionServer(&server{})
}

func (self *server) Login(ctx proto.Context, request *proto.LoginRequest) (*proto.LoginReply, error) {
	tk := token.Entoken(request.Name)
	reply := &proto.LoginReply{}
	reply.Flag = 0
	reply.Token = tk
	reply.Lease = 600
	reply.Interval = 500

	return reply, nil
}

func (self *server) Logout(ctx proto.Context, request *proto.LogoutRequest) (*proto.LogoutReply, error) {
	reply := &proto.LogoutReply{}
	return reply, nil
}
