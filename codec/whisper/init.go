package whisper

import (
	"github.com/hitzhangjie/go-rpc/codec"
)

const (
	WhisperClientCodec = "whisper_client_codec"
	WhisperServerCodec = "whisper_server_codec"
)

func init() {
	codec.Mux.Lock()

	c := &ClientCodec{}
	s := &ServerCodec{}
	codec.CodecMappings[WhisperClientCodec] = c
	codec.CodecMappings[WhisperServerCodec] = s
	codec.ReaderMappings[WhisperClientCodec] = codec.NewMessageReader(c)
	codec.ReaderMappings[WhisperServerCodec] = codec.NewMessageReader(s)
	codec.Mux.Unlock()
}
