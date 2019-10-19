package main

import (
	"github.com/hitzhangjie/go-rpc/codec/whisper"
	"github.com/hitzhangjie/go-rpc/server"
)

func main() {
	tcpSvr, err := server.NewTcpServer("tcp4", "127.0.0.1:8888", whisper.Whisper)
	if err != nil {
		panic(err)
	}
	tcpSvr.Start()
}
