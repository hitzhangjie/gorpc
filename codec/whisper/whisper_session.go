package whisper

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hitzhangjie/gorpc-framework/codec"
)

type WhisperSession struct {
	codec.BaseSession
}

func (s *WhisperSession) RPCName() string {
	return s.Request().(*Request).GetRpcname()
}

func (s *WhisperSession) SetErrorResponse(err error) {
	rsp := s.Response().(*Response)
	rsp.ErrCode = proto.Uint32(10000)
	rsp.ErrMsg = proto.String(err.Error())
}

func (s *WhisperSession) TraceContext() interface{} {
	req := s.Request().(*Request)
	if req.Meta != nil {
		return []byte(req.Meta["traceContext"])
	}
	return nil
}

// WhisperSessionBuilder builder for WhisperSession
type WhisperSessionBuilder struct{}

// Build build a new WhisperSession
func (b *WhisperSessionBuilder) Build(req interface{}) (codec.Session, error) {
	return newSession(req)
}

func newSession(req interface{}) (codec.Session, error) {

	reqHead, ok := req.(*Request)
	if !ok {
		return nil, fmt.Errorf("req:%v not *whisper.ReqHead", req)
	}

	rspHead := &Response{}
	rspHead.Seqno = reqHead.Seqno

	session := &WhisperSession{
		codec.BaseSession{
			ReqHead: reqHead,
			RspHead: rspHead,
		},
	}

	return session, nil
}
