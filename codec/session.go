package codec

type Session interface {
	RPC() string

	Request() interface{}
	SetRequest(req interface{})

	Response() interface{}
	SetResponse(rsp interface{})

	TraceContext() interface{}
	//TraceStart() func(Session)
	//TraceFinish() func(Session)
}
