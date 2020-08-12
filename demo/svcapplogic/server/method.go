package service

import (
	"iogo/demo/proto/applogic"
)

type server struct {
}

func init() {
	proto.RegisterApplogicServer(&server{})
}

func (self *server) Say(context proto.Context, request *proto.SayRequest) (*proto.SayReply, error) {
	Reply := &proto.SayReply{}
	Reply.Flag = "ok"
	return Reply, nil
}
