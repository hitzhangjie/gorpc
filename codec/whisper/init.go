package whisper

import (
	"github.com/hitzhangjie/gorpc-framework/codec"
)

const Whisper = "whisper"

func init() {
	codec.RegisterCodec(Whisper, &ServerCodec{}, &ClientCodec{})
	codec.RegisterSessionBuilder(Whisper, &WhisperSessionBuilder{})
}
