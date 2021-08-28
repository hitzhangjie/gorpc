package codec

import (
	"bytes"
	"compress/gzip"
	"io"
)

// Compressor compress and decompress the data
type Compressor interface {
	// Compress return compressed data, return error if any error encountered
	Compress(data []byte) ([]byte, error)

	// Decompress return decompressed data, return error if any error encountered
	Decompress([]byte) ([]byte, error)
}

// CompressType compress type
type CompressType int

const (
	CompressGZip   = CompressType(iota) // gzip compress
	CompressSnappy                      // snappy compress
	CompressLZ4                         // lz4 compress
)

// GZipCompressor compressor using GZip cmopression algorithm
type GZipCompressor struct{}

// Compress returns the compressed format of data
func (c *GZipCompressor) Compress(data []byte) ([]byte, error) {
	buf := &bytes.Buffer{}
	w, err := gzip.NewWriterLevel(buf, gzip.BestSpeed)
	if err != nil {
		return nil, err
	}

	_, err = w.Write(data)
	if err != nil {
		return nil, err
	}
	w.Close()

	return buf.Bytes(), nil
}

// Decompress decompress returns the compressed form of data
func (c *GZipCompressor) Decompress(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return io.ReadAll(r)
}
