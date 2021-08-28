package codec_test

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"

	"github.com/hitzhangjie/gorpc/codec"
	pb "github.com/hitzhangjie/gorpc/testdata/protobuf/helloworld"
)

var marshaler = &codec.PBMarshaler{}

func TestPBMarshaler_Marshal(t *testing.T) {

	t.Run("invalid pkg", func(t *testing.T) {
		req := pb.HelloRequest{
			Msg: "helloworld",
		}
		b, err := marshaler.Marshal(req)
		assert.Equal(t, codec.ErrMarshalInvalidPB, err)
		assert.Nil(t, b)
	})

	t.Run("valid pkg", func(t *testing.T) {
		req := &pb.HelloRequest{
			Msg: "hellowrold",
		}
		b, err := marshaler.Marshal(req)
		assert.Nil(t, err)

		b2, err := proto.Marshal(req)
		if err != nil {
			panic(err)
		}
		assert.Equal(t, b2, b)
	})
}

func TestPBMarshaler_Unmarshal(t *testing.T) {

	req := pb.HelloRequest{
		Msg: "helloworld",
	}
	b, err := proto.Marshal(&req)
	if err != nil {
		panic(err)
	}

	t.Run("invalid pkg", func(t *testing.T) {
		r := pb.HelloRequest{}
		err := marshaler.Unmarshal(b, r)
		assert.Equal(t, codec.ErrMarshalInvalidPB, err)
	})

	t.Run("valid pkg", func(t *testing.T) {
		r := pb.HelloRequest{}
		err := marshaler.Unmarshal(b, &r)
		assert.Nil(t, err)
		assert.Equal(t, req.Msg, r.Msg)
	})
}
