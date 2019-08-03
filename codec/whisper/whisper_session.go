package whisper

import (
	"github.com/golang/protobuf/proto"
	"github.com/hitzhangjie/go-rpc/codec"
)

type WhisperSession struct {
	codec.BaseSession
}

func (w *WhisperSession) RPC() string {
	panic("implement me")
}

func (w *WhisperSession) Request() interface{} {
	panic("implement me")
}

func (w *WhisperSession) SetRequest(req interface{}) {
	panic("implement me")
}

func (w *WhisperSession) Response() interface{} {
	panic("implement me")
}

func (w *WhisperSession) SetResponse(rsp interface{}) {
	panic("implement me")
}

func (w *WhisperSession) SetErrorResponse(err error) {
	panic("implement me")
}

func (w *WhisperSession) TraceContext() interface{} {
	panic("implement me")
}

func NewSession(req []byte) (codec.Session, error) {

	reqHead := &Request{}
	if err := proto.Unmarshal(req, reqHead); err != nil {
		return nil, err
	}
	rspHead := &Response{}
	rspHead.Seqno = reqHead.Seqno

	session := &WhisperSession{
		codec.BaseSession{
			Request:  reqHead,
			Response: rspHead,
		},
	}

	return session, nil
}
