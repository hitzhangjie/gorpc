package codec_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hitzhangjie/gorpc/codec"
)

var (
	msg = "helloworld"

	compressor = &codec.GZipCompressor{}
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
