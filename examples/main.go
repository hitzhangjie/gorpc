package main

import (
	"github.com/hitzhangjie/go-rpc/codec/whisper"
	"github.com/hitzhangjie/go-rpc/server"
)

func main() {
	svr, err := server.NewServer()
	if err != nil {
		panic(err)
	}

	tcpSvr, err := server.NewTcpServer("tcp4", "127.0.0.1:8888", whisper.WhisperServerCodec)
	if err != nil {
		panic(err)
	}
	tcpSvr.Register(svr)

	udpSvr, err := server.NewUdpServer("udp4", "127.0.0.1:8888", whisper.WhisperServerCodec)
	if err != nil {
		panic(err)
	}
	udpSvr.Register(svr)

	svr.Start()
}
