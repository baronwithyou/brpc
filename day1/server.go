package day1

import (
	"brpc/day1/codec"
	"encoding/json"
	"io"
	"log"
	"net"
	"reflect"
)

const MagicNumber = 0x3bef5c

type Option struct {
	MagicNumber int        // MagicNumber marks this's a geerpc request
	CodecType   codec.Type // client may choose different Codec to encode body
}

var DefaultOption = &Option{
	MagicNumber: MagicNumber,
	CodecType:   codec.GobType,
}

// Server represents an RPC Server.
type Server struct{}

// NewServer returns a new Server.
func NewServer() *Server {
	return &Server{}
}

// DefaultServer is the default instance of *Server.
var DefaultServer = NewServer()

// Accept accepts connections on the listener and serves requests
// for each incoming connection.
func (server *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("rpc server: accept error:", err)
			return
		}
		go server.ServeConn(conn)
	}
}

// Accept accepts connections on the listener and serves requests
// for each incoming connection.
func Accept(lis net.Listener) { DefaultServer.Accept(lis) }

func (server *Server) ServeConn(conn io.ReadWriteCloser) {
	var option Option
	// 处理option
	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&option); err != nil {
		log.Println("rpc decode option err :", err)
		return
	}

	if option.MagicNumber != MagicNumber {
		log.Println("rpc magic number is wrong")
		return
	}

	f := codec.NewCodecFuncMap[option.CodecType]

	server.serveCodec(f(conn))
}

func (server *Server) serveCodec(cc codec.Codec) {

	for {
		// readRequest
		req, err := server.readRequest(cc)
		if err != nil {
			// sendResponse
			server.sendResponse()
		}

		// handleRequest
		go server.handleRequest()
	}

}

func (server *Server) readRequest(cc codec.Codec) (*request, error) {
	return nil, nil
}

func (server *Server) sendResponse(cc codec.Codec, h *codec.Header, body interface{}) {

}

func (server *Server) handleRequest() {

}

type request struct {
	h            *codec.Header // header of request
	argv, replyv reflect.Value // argv and replyv of request
}
