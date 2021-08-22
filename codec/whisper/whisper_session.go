package whisper

import (
	"fmt"

	"github.com/golang/protobuf/proto"

	"github.com/hitzhangjie/gorpc/codec"
	"github.com/hitzhangjie/gorpc/errors"
)

type WhisperSession struct {
	codec.BaseSession
}

func (s *WhisperSession) RPCName() string {
	return s.Request().(*Request).GetRpcname()
}

func (s *WhisperSession) SetError(err error) {
	var code int
	var msg string
	rsp := s.Response().(*Response)
	if errors.IsFrameworkError(err) {
		code = errors.ErrorCode(err)
		msg = errors.ErrorMsg(err)
	} else {
		code = 10000
		msg = err.Error()
	}
	rsp.ErrCode = proto.Uint32(uint32(code))
	rsp.ErrMsg = proto.String(msg)
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
