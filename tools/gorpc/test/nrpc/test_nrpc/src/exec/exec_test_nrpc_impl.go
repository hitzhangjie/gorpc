package exec

import (
    "context"
	"git.code.oa.com/nrpc_protos/test_nrpc"
	"git.code.oa.com/go-neat/core/nserver/nsession"
	
)

func BuyAppleImpl(ctx context.Context, session nsession.NSession, req *test_nrpc.BuyAppleReq, rsp *test_nrpc.BuyAppleRsp) error {
	// business logic
	return nil
}

func SellAppleImpl(ctx context.Context, session nsession.NSession, req *test_nrpc.SellAppleReq, rsp *test_nrpc.SellAppleRsp) error {
	// business logic
	return nil
}

