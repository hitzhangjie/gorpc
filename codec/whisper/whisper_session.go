package whisper

import (
	"fmt"
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

func NewSession(req interface{}) (codec.Session, error) {

	reqHead, ok := req.(*Request)
	if !ok {
		return nil, fmt.Errorf("req:%v not *whisper.Request", req)
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
