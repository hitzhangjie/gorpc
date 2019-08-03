package whisper

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/gogo/protobuf/proto"
	"github.com/hitzhangjie/go-rpc/codec"
)

const maxWhisperPkgSize = 64 * (2 << 10) // 64KB

// ServerCodec serverside codec
type ServerCodec struct{}

func (s *ServerCodec) Name() string {
	return Whisper
}

func (s *ServerCodec) Encode(pkg interface{}) ([]byte, error) {

	pb, ok := pkg.(*Response)
	if !ok {
		return nil, errors.New("pkg not valid *whisper.Response")
	}

	data, err := proto.Marshal(pb)
	if err != nil {
		return nil, err
	}

	b := &bytes.Buffer{}
	binary.Write(b, binary.BigEndian, int8(0x38))
	binary.Write(b, binary.BigEndian, int32(len(data)))
	binary.Write(b, binary.BigEndian, data)
	binary.Write(b, binary.BigEndian, int8(0x49))

	return b.Bytes(), nil
}

func (s *ServerCodec) Decode(in []byte, message interface{}) error {

	request, ok := message.(*Request)
	if !ok {
		return errors.New("message not *whisper.Request")
	}

	if len(in) < 5 {
		return codec.CodecReadIncomplete
	}

	b := bytes.NewBuffer(in)

	// pkg: | 1B:0x38 | 4B:len | payload | 1B: 0x49 |
	var (
		pkgStx int8
		pkgLen int32
		pkgEtx int8
		err    error
	)
	// stx
	if err = binary.Read(b, binary.BigEndian, &pkgStx); err != nil {
		return err
	}
	if pkgStx != 0x38 {
		return codec.CodecReadInvalid
	}
	// len
	if err = binary.Read(b, binary.BigEndian, &pkgLen); err != nil {
		return err
	}
	if pkgLen > maxWhisperPkgSize {
		return codec.CodecReadTooBig
	}
	if len(in) != int(1+4+pkgLen+1) {
		return codec.CodecReadIncomplete
	}
	// etx
	if err = binary.Read(b, binary.BigEndian, &pkgEtx); err != nil {
		return err
	}
	if pkgEtx != 0x49 {
		return codec.CodecReadInvalid
	}
	// payload
	payload := &bytes.Buffer{}
	if err = binary.Read(b, binary.BigEndian, payload); err != nil {
		return codec.CodecReadError
	}
	if err = proto.Unmarshal(payload.Bytes(), request); err != nil {
		return err
	}

	return nil
}

// ClientCodec clientside codec
type ClientCodec struct{}

func (c *ClientCodec) Name() string {
	return Whisper
}

func (c *ClientCodec) Encode(pkg interface{}) ([]byte, error) {
	return nil, nil
}

func (c *ClientCodec) Decode([]byte, interface{}) error {
	return nil
}

func (c *ClientCodec) Session([]byte) (codec.Session, error) {
	return nil, nil
}
