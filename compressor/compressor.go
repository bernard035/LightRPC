package compressor

type CompressType uint16

const (
	Raw CompressType = iota
	Gzip
	Zlib
	Snappy
)

// Compressors which supported by rpc
var Compressors = map[CompressType]Compressor{
	Raw:    RawCompressor{},
	Gzip:   GzipCompressor{},
	Zlib:   ZlibCompressor{},
	Snappy: SnappyCompressor{},
}

// Compressor is interface, each compressor has Zip and Unzip functions
type Compressor interface {
	Zip([]byte) ([]byte, error)
	Unzip([]byte) ([]byte, error)
}
