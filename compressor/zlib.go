package compressor

import (
	"bytes"
	"compress/zlib"
	"io"
)

// ZlibCompressor implements the Compressor interface
type ZlibCompressor struct {
}

// Zip .
func (_ ZlibCompressor) Zip(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	writer := zlib.NewWriter(buf)
	defer writer.Close()
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
func (_ ZlibCompressor) Unzip(data []byte) ([]byte, error) {
	reader, err := zlib.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	data, err = io.ReadAll(reader)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return nil, err
	}
	return data, nil
}
