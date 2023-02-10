package compressor

import (
	"bytes"
	"github.com/golang/snappy"
	"io"
)

// SnappyCompressor implements the Compressor interface
type SnappyCompressor struct {
}

// Zip .
func (_ SnappyCompressor) Zip(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	writer := snappy.NewBufferedWriter(buf)
	defer func() {
		writer.Close()
	}()
	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}
	err = writer.Flush()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

// Unzip .
func (_ SnappyCompressor) Unzip(data []byte) ([]byte, error) {
	reader := snappy.NewReader(bytes.NewBuffer(data))
	data, err := io.ReadAll(reader)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return nil, err
	}
	return data, nil
}
