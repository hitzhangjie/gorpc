package main

import (
	"context"

	"github.com/hitzhangjie/gorpc/codec/whisper"
	"github.com/hitzhangjie/gorpc/server"
)

func main() {
	service := server.NewService("testsvr")

	ctx := context.Background()

	if err := service.ListenAndServe(ctx, "tcp4", "127.0.0.1:8888", whisper.Whisper); err != nil {
		panic(err)
	}

	if err := service.ListenAndServe(ctx, "udp4", "127.0.0.1:8888", whisper.Whisper); err != nil {
		panic(err)
	}
}
