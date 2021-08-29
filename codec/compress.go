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

// gzipCompressor compressor using GZip cmopression algorithm
type gzipCompressor struct {
	reader *gzip.Reader
	writer *gzip.Writer
	buffer *bytes.Buffer
}

// NewGZipCompressor returns a new initialized gzipCompressor
func NewGZipCompressor() *gzipCompressor {
	return &gzipCompressor{
		reader: new(gzip.Reader),
		writer: new(gzip.Writer),
		buffer: bytes.NewBuffer(nil),
	}
}

// Compress returns the compressed format of data
func (c *gzipCompressor) Compress(data []byte) ([]byte, error) {
	buf := &bytes.Buffer{}

	c.writer.Reset(buf)
	if _, err := c.writer.Write(data); err != nil {
		return nil, err
	}
	c.writer.Close()

	return buf.Bytes(), nil
}

// Decompress decompress returns the compressed form of data
func (c *gzipCompressor) Decompress(data []byte) ([]byte, error) {
	if err := c.reader.Reset(bytes.NewReader(data)); err != nil {
		return nil, err
	}
	defer c.reader.Close()

	return io.ReadAll(c.reader)
}
