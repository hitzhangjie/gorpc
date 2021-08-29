package codec_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hitzhangjie/gorpc/codec"
)

type fakeCodec struct{}

func (c *fakeCodec) Name() string {
	return "fake"
}

func (c *fakeCodec) Encode(pkg interface{}) (dat []byte, err error) {
	return nil, nil
}

func (c *fakeCodec) Decode(dat []byte) (req interface{}, n int, err error) {
	return nil, 0, nil
}

func TestRegisterCodec(t *testing.T) {
	c := &fakeCodec{}
	codec.RegisterCodec("fake", c, c)
	assert.Equal(t, codec.ServerCodec("fake"), c)
	assert.Equal(t, codec.ClientCodec("fake"), c)
}
