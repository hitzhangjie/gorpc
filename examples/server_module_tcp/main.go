package main

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/hitzhangjie/gorpc/client"
	"github.com/hitzhangjie/gorpc/client/selector"
	"github.com/hitzhangjie/gorpc/codec/whisper"
	"github.com/hitzhangjie/gorpc/router"
	"github.com/hitzhangjie/gorpc/transport"
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

		tcpSvr, err := transport.NewTcpServerTransport(context.TODO(),
			"tcp4", "127.0.0.1:8888", whisper.Whisper, transport.WithRouter(r))
		if err != nil {
			panic(err)
		}
		tcpSvr.ListenAndServe()
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
		fmt.Println("client recv response:", rspHead)
	}()

	select {}
}
