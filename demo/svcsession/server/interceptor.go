package service

import (
	"fmt"

	"iogo"
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
	isTrue := false
	flag := fmt.Sprintf("%p", info)
	if isTrue {
		inh, _ := md.FromIncomingContext(ctx)
		iogo.Log.I(LC_SVCSESSION_CALLBEFOR, "UnaryIntercaptorBefor:%s call: %s request: %+v, in-header: %+v", flag, info.FullMethod, req, inh)
	}
	resp, err = handler(ctx, req)
	if isTrue {
		if err != nil {
			iogo.Log.I(LC_SVCSESSION_CALLAFTER, "UnaryIntercaptorAfter:%s call: %s error: %v", flag, info.FullMethod, err)
		} else {
			iogo.Log.I(LC_SVCSESSION_TOKENINFO, "UnaryIntercaptorAfter:%s call: %s reply: %+v", flag, info.FullMethod, resp)
		}
	}
	return
}
