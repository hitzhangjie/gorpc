package main

import (
	"github.com/hitzhangjie/go-rpc/codec/whisper"
	"github.com/hitzhangjie/go-rpc/server"
)

func main() {
	svr := server.NewService()

	tcpSvr, err := server.NewTcpServerTransport("tcp4", "127.0.0.1:8888", whisper.Whisper)
	if err != nil {
		panic(err)
	}
	tcpSvr.Register(svr)

	udpSvr, err := server.NewUdpServerTransport("udp4", "127.0.0.1:8888", whisper.Whisper)
	if err != nil {
		panic(err)
	}
	udpSvr.Register(svr)

	svr.Start()
}
