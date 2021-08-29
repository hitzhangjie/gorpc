package codec_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hitzhangjie/gorpc/codec"
)

var (
	msg = "helloworld"

	compressor = codec.NewGZipCompressor()
)

func TestGZipCompressor_Compress(t *testing.T) {
	b, err := compressor.Compress([]byte(msg))
	assert.Nil(t, err)
	t.Logf("compressed data hex: %s, len: %d", hex.EncodeToString(b), len(b))
}

func TestGZipCompressor_Decompress(t *testing.T) {
	b, err := compressor.Compress([]byte(msg))
	assert.Nil(t, err)
	assert.Len(t, b, 38)

	dat, err := compressor.Decompress(b)
	assert.Nil(t, err)
	assert.Equal(t, msg, string(dat))
}

func BenchmarkGzipCompressor_Decompress(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.Run("", func(b *testing.B) {
			buf, err := compressor.Compress([]byte(msg))
			assert.Nil(b, err)

			dat, err := compressor.Decompress(buf)
			assert.Nil(b, err)
			assert.Equal(b, msg, string(dat))
		})
	}
}
