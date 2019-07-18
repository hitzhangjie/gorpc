package exec

import (
	"git.code.oa.com/go-neat/core/nserver/default_nserver"
)

func init() {
	//注册服务接口
    default_nserver.AddExec("BuyApple", BuyApple)
    default_nserver.AddExec("SellApple", SellApple)
}
