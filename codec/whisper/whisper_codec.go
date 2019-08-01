package whisper

// ClientCodec clientside codec
type ClientCodec struct{}

func (c *ClientCodec) Name() string {
	return WhisperClientCodec
}

func (c *ClientCodec) Encode(pkg interface{}) ([]byte, error) {
	return nil, nil
}

func (c *ClientCodec) Decode([]byte, interface{}) error {
	return nil
}

// ServerCodec serverside codec
type ServerCodec struct{}

func (s *ServerCodec) Name() string {
	return WhisperServerCodec
}

func (s *ServerCodec) Encode(pkg interface{}) ([]byte, error) {
	return nil, nil
}

func (s *ServerCodec) Decode([]byte, interface{}) error {
	return nil
}
