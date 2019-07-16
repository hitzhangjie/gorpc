package server

type Server interface {
	Start()
	Stop()
	Handle(h Handler) error
}

type server struct {
	mods [] *serverModule
}

type serverModule interface {
	Start()
	Stop()
}

type StreamServer struct {
	svr *server
}

func (s *StreamServer) Start() {

}

func (s *StreamServer) Stop() {

}

type PacketServer struct {
	svr *server
}

func (s *PacketServer) Start() {
}

func (s *PacketServer) Stop() {

}
