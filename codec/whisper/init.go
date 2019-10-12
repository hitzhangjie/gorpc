package whisper

import (
	"github.com/hitzhangjie/go-rpc/codec"
)

const Whisper = "whisper"

func init() {
	codec.RegisterCodec(Whisper, &ServerCodec{}, &ClientCodec{})
	codec.RegisterSessionBuilder(Whisper, &WhisperSessionBuilder{})
}
