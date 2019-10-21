package whisper

import (
	"github.com/golang/protobuf/proto"
	"testing"
)

var serverCodec = &ServerCodec{}
var clientCodec = &ClientCodec{}

// test server
func Test_ServerEncode_ClientDecode(t *testing.T) {
	response := &Response{
		Seqno:   proto.Uint64(12345),
		ErrCode: proto.Uint32(0),
		ErrMsg:  proto.String("Success"),
		Body:    nil,
	}

	dat, err := serverCodec.Encode(response)
	if err != nil {
		t.Fatalf("serverCodec encode, err:%v", err)
	}
	t.Logf("serverCodec encode, ok, data:%v", dat)

	rsp, n, err := clientCodec.Decode(dat)
	if err != nil {
		t.Fatalf("clientCodec decode, err:%v", err)
	}
	t.Logf("clientCodec decode, ok, data:%v", rsp)

	if n != len(dat) {
		t.Fatalf("clientCodec data len:%v, want:%v", n, len(dat))
	}

	// protobuf message has some internal fields, it will change, so we cannot use reflect.DeepEqual(response, drsp).
	drsp := rsp.(*Response)
	if *response.Seqno != *drsp.Seqno ||
		*response.ErrCode != *drsp.ErrCode ||
		*response.ErrMsg != *drsp.ErrMsg ||
		len(drsp.Body) != 0 {
		t.Fatalf("clientCodec rsp:%v, want:%v", rsp, response)
	}
}

// test client
func Test_ClientCodec_ServerDecode(t *testing.T) {
	request := &Request{
		Seqno:   proto.Uint64(1),
		Appid:   proto.String("APPID_100"),
		Userid:  proto.String("1194606858"),
		Userkey: proto.String("hello"),
		Version: proto.Uint32(1),
	}

	dat, err := clientCodec.Encode(request)
	if err != nil {
		t.Fatalf("clientCodec encode, err:%v, data:%v", err, dat)
	}
	t.Logf("clientCodec, ok, data:%v", dat)

	req, _, err := serverCodec.Decode(dat)
	if err != nil {
		t.Fatalf("serverCodec decode, err:%v, data:%v", err, req)
	}
	t.Logf("serverCodec, ok, data:%v", req)

	// protobuf message has some internal fields, it will change, so we cannot use reflect.DeepEqual(response, drsp).
	dreq := req.(*Request)
	if *request.Seqno != *dreq.Seqno ||
		*request.Appid != *dreq.Appid ||
		*request.Userid != *dreq.Userid ||
		*request.Userkey != *dreq.Userkey ||
		*request.Version != *dreq.Version {
		t.Fatalf("serverCodec req:%v, want:%v", dreq, request)
	}
}
