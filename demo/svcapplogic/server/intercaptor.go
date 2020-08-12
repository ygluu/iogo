package service

import (
	"errors"
	"fmt"

	"iogo"
	"iogo/demo/comm/token"
	"iogo/unit/md"
)

func init() {
	iogo.RegisterIntercaptors(unaryIntercaptor, streamInterceptor)
}

func streamInterceptor(srv interface{}, ss interface{}, info *iogo.StreamServerInfo, handler iogo.StreamHandler) error {
	return handler(srv, ss)
}

func unaryIntercaptor(ctx iogo.Context, req interface{}, info *iogo.UnaryServerInfo,
	handler iogo.UnaryHandler) (resp interface{}, err error) {
	inh, _ := md.FromIncomingContext(ctx)
	if inh == nil {
		err = errors.New("Invalid access")
		return
	}
	tks := inh["set-token"]
	if tks == nil || len(tks) == 0 {
		err = errors.New("There is no token")
		return
	}
	if !token.CheckToken(tks[0]) {
		err = errors.New("Invalid token")
		return
	}

	isTrue := false
	flag := fmt.Sprintf("%p", info)
	if isTrue {
		iogo.Log.I(LC_SVCAPPL_CALLBEFOR, "UnaryIntercaptorBefor:%s call: %s request: %+v, in-header: %+v", flag, info.FullMethod, req, inh)
	}
	resp, err = handler(ctx, req)
	if isTrue {
		if err != nil {
			iogo.Log.I(LC_SVCAPPL_CALLAFTER, "UnaryIntercaptorAfter:%s call: %s error: %v", flag, info.FullMethod, err)
		} else {
			iogo.Log.I(LC_SVCAPPL_CALLAFTER, "UnaryIntercaptorAfter:%s call: %s reply: %+v", flag, info.FullMethod, resp)
		}
	}
	return
}
