package whisper

type WhisperSession struct{}

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

func (w *WhisperSession) TraceContext() interface{} {
	panic("implement me")
}


