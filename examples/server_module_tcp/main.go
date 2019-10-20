package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hitzhangjie/go-rpc/client"
	"github.com/hitzhangjie/go-rpc/client/selector"
	"github.com/hitzhangjie/go-rpc/codec/whisper"
	"github.com/hitzhangjie/go-rpc/router"
	"github.com/hitzhangjie/go-rpc/server"
	"time"
)

func main() {
	// server
	go func() {
		r := router.NewRouter()
		r.Forward("/hello", func(ctx context.Context, req interface{}) (rsp interface{}, err error) {
			pbreq := req.(*whisper.Request)
			fmt.Printf("server recv req:%v", pbreq)

			pbrsp := &whisper.Response{
				Seqno:   pbreq.Seqno,
				ErrCode: proto.Uint32(0),
				ErrMsg:  proto.String("success"),
			}
			return pbrsp, nil
		})

		tcpSvr, err := server.NewTcpServer("tcp4", "127.0.0.1:8888", whisper.Whisper, server.WithRouter(r))
		if err != nil {
			panic(err)
		}
		tcpSvr.Start()
	}()

	// client
	go func() {
		time.Sleep(time.Second)
		cli := client.NewClient("test",
			client.WithAddress("127.0.0.1:8888"),
			client.WithCodec("whisper"),
			client.WithTransportType(client.TCP4),
			client.WithSelector(selector.NewIPSelector("tcp4", []string{"127.0.0.1:8888"})),
		)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		reqHead := &whisper.Request{
			Seqno:   proto.Uint64(1000),
			Appid:   proto.String("strong"),
			Rpcname: proto.String("/hello"),
			Userid:  proto.String("zhangjie"),
			Userkey: proto.String("can"),
			Version: proto.Uint32(0),
		}
		fmt.Println(reqHead)

		v, err := cli.Invoke(ctx, reqHead)
		if err != nil {
			fmt.Println(err)
			return
		}

		rspHead := v.(*whisper.Response)
		fmt.Println(rspHead)
	}()

	select {}
}
