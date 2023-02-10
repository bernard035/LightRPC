package LightRPC

import (
	"LightRPC/codec"
	"LightRPC/serializer"
	"log"
	"net"
	"net/rpc"
)

// Server rpc server based on net/rpc implementation
type Server struct {
	*rpc.Server
	serializer.Serializer
}

// NewServer Create a new rpc server
func NewServer(opts ...Option) *Server {
	options := options{
		serializer: serializer.Proto,
	}
	for _, option := range opts {
		option(&options)
	}

	return &Server{&rpc.Server{}, options.serializer}
}

// Register register rpc function
func (s *Server) Register(receiver interface{}) error {
	return s.Server.Register(receiver)
}

// RegisterName register the rpc function with the specified name
func (s *Server) RegisterName(name string, receiver interface{}) error {
	return s.Server.RegisterName(name, receiver)
}

// Serve start service
func (s *Server) Serve(lis net.Listener) {
	log.Printf("tinyrpc started on: %s", lis.Addr().String())
	for {
		conn, err := lis.Accept()
		if err != nil {
			continue
		}
		go s.Server.ServeCodec(codec.NewServerCodec(conn, s.Serializer))
	}
}
