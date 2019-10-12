package whisper

import (
	"github.com/golang/protobuf/proto"
	"testing"
)

var serverCodec = &ServerCodec{}
var clientCodec = &ClientCodec{}

// test server
func TestServerCodec(t *testing.T) {
	response := &Response{
		Seqno: proto.Uint64(1),
	}

	dat, err := serverCodec.Encode(response)
	t.Logf("serverCodec, err:%v, marshal:%v", err, dat)

	rsp, err := serverCodec.Decode(dat)
	t.Logf("serverCodec, err:%v, unmarshal:%v", err, rsp)
}

// test client
func TestClientCodec(t *testing.T) {
	request := &Request{
		Seqno: proto.Uint64(1),
	}

	dat, err := clientCodec.Encode(request)
	t.Logf("clientCodec, err:%v, marshal:%v", err, dat)

	rsp, err := clientCodec.Decode(dat)
	t.Logf("clientCodec, err:%v, unmarshal:%v", err, rsp)
}
