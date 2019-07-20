package exec

import (
	"context"
	"git.code.oa.com/go-neat/core/nserver"
	"git.code.oa.com/go-neat/core/nserver/nsession"
	"git.code.oa.com/go-neat/tencent/attr"
	"git.code.oa.com/nrpc_protos/test_nrpc"
)

func BuyApple(ctx context.Context, session nsession.NSession) (interface{}, error) {
	req := &test_nrpc.BuyAppleReq{}
	err := session.ParseRequestBody(req)
	
	if err != nil {
		attr.Monitor(0, 1) //test_nrpc.BuyAppleReq解析失败
		session.Logger().Error("parse req err %v", err)
		return nil, err
	}
	
	rsp := &test_nrpc.BuyAppleRsp{}
	err = BuyAppleImpl(ctx, session, req, rsp)
	if err != nil {
		attr.Monitor(0, 1) //test_nrpc.BuyAppleReq处理异常
		session.Logger().Error("handle req err %v", err)

		if _, ok := err.(nserver.Error); ok {
        	return nil, err
        }
        return nil, nserver.CreateError(nserver.EXEC_HANDLE_ERROR, err)
	}
	
	return rsp, nil
}

func SellApple(ctx context.Context, session nsession.NSession) (interface{}, error) {
	req := &test_nrpc.SellAppleReq{}
	err := session.ParseRequestBody(req)
	
	if err != nil {
		attr.Monitor(0, 1) //test_nrpc.SellAppleReq解析失败
		session.Logger().Error("parse req err %v", err)
		return nil, err
	}
	
	rsp := &test_nrpc.SellAppleRsp{}
	err = SellAppleImpl(ctx, session, req, rsp)
	if err != nil {
		attr.Monitor(0, 1) //test_nrpc.SellAppleReq处理异常
		session.Logger().Error("handle req err %v", err)

		if _, ok := err.(nserver.Error); ok {
        	return nil, err
        }
        return nil, nserver.CreateError(nserver.EXEC_HANDLE_ERROR, err)
	}
	
	return rsp, nil
}

