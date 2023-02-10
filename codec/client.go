package codec

import (
	"LightRPC/compressor"
	"LightRPC/header"
	"LightRPC/serializer"
	"bufio"
	"hash/crc32"
	"io"
	"net/rpc"
	"sync"
)

type clientCodec struct {
	reader io.Reader
	writer io.Writer
	closer io.Closer

	compressor compressor.CompressType // rpc compress type(raw,gzip,zlib,snappy)
	serializer serializer.Serializer
	response   header.ResponseHeader // rpc response header
	mutex      sync.Mutex            // protect pending map
	pending    map[uint64]string
}

// NewClientCodec Create a new client codec
func NewClientCodec(conn io.ReadWriteCloser,
	compressType compressor.CompressType, serializer serializer.Serializer) rpc.ClientCodec {

	return &clientCodec{
		reader:     bufio.NewReader(conn),
		writer:     bufio.NewWriter(conn),
		closer:     conn,
		compressor: compressType,
		serializer: serializer,
		pending:    make(map[uint64]string),
	}
}

// WriteRequest Write the rpc request header and body to the io stream
func (c *clientCodec) WriteRequest(r *rpc.Request, param interface{}) error {
	c.mutex.Lock()
	c.pending[r.Seq] = r.ServiceMethod
	c.mutex.Unlock()

	if _, ok := compressor.Compressors[c.compressor]; !ok {
		return NotFoundCompressorError
	}
	reqBody, err := c.serializer.Marshal(param)
	if err != nil {
		return err
	}
	compressedReqBody, err := compressor.Compressors[c.compressor].Zip(reqBody)
	if err != nil {
		return err
	}
	h := header.RequestPool.Get().(*header.RequestHeader)
	defer func() {
		h.ResetHeader()
		header.RequestPool.Put(h)
	}()
	h.ID = r.Seq
	h.Method = r.ServiceMethod
	h.RequestLen = uint32(len(compressedReqBody))
	h.CompressType = compressor.CompressType(c.compressor)
	h.Checksum = crc32.ChecksumIEEE(compressedReqBody)

	if err := sendFrame(c.writer, h.Marshal()); err != nil {
		return err
	}
	if err := write(c.writer, compressedReqBody); err != nil {
		return err
	}

	c.writer.(*bufio.Writer).Flush()
	return nil
}

// ReadResponseHeader read the rpc response header from the io stream
func (c *clientCodec) ReadResponseHeader(r *rpc.Response) error {
	c.response.ResetHeader()
	data, err := recvFrame(c.reader)
	if err != nil {
		return err
	}
	err = c.response.Unmarshal(data)
	if err != nil {
		return err
	}
	c.mutex.Lock()
	r.Seq = c.response.ID
	r.Error = c.response.Error
	r.ServiceMethod = c.pending[r.Seq]
	delete(c.pending, r.Seq)
	c.mutex.Unlock()
	return nil
}

// ReadResponseBody read the rpc response body from the io stream
func (c *clientCodec) ReadResponseBody(param interface{}) error {
	if param == nil {
		if c.response.ResponseLen != 0 {
			if err := read(c.reader, make([]byte, c.response.ResponseLen)); err != nil {
				return err
			}
		}
		return nil
	}

	respBody := make([]byte, c.response.ResponseLen)
	err := read(c.reader, respBody)
	if err != nil {
		return err
	}

	if c.response.Checksum != 0 {
		if crc32.ChecksumIEEE(respBody) != c.response.Checksum {
			return UnexpectedChecksumError
		}
	}

	if c.response.GetCompressType() != c.compressor {
		return CompressorTypeMismatchError
	}

	resp, err := compressor.Compressors[c.response.GetCompressType()].Unzip(respBody)
	if err != nil {
		return err
	}

	return c.serializer.Unmarshal(resp, param)
}

func (c *clientCodec) Close() error {
	return c.closer.Close()
}
