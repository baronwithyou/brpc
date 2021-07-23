package day1

import (
	"brpc/day1/codec"
	"encoding/json"
	"io"
	"log"
	"net"
	"reflect"
	"sync"
)

/**
* 几个重点：
* 1. json.Decoder.decode()的时候，允许有其他类型的数据。
* 2. 锁：数据发送的时候需要上锁，否则会有杂糅的可能性（同时还得考虑同步+异步的场景，Line93/Line98）。
* 3. reflect
 */

type request struct {
	h            *codec.Header // header of request
	argv, replyv reflect.Value // argv and replyv of request
}

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
		// 每个connection使用一个goroutine来handle
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

// invalidRequest is a placeholder for response argv when error occurs
var invalidRequest = struct{}{}

func (server *Server) serveCodec(cc codec.Codec) {
	sending := new(sync.Mutex)
	wg := new(sync.WaitGroup)

	// 多个请求体会出现在一个request里面，所以需要轮询
	for {
		// readRequest
		req, err := server.readRequest(cc)
		if err != nil {
			// sendResponse
			// 因为handleRequest是异步的，跟这个方法会有同步发出的可能，所以得上锁。
			server.sendResponse(cc, req.h, invalidRequest, sending)
		}

		// handleRequest - 如果不使用goroutine，请求会串行执行。
		// 需要使用锁，否则消息会杂糅在一起
		go server.handleRequest(cc, req, sending, wg)
	}
}

func (server *Server) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
	var h codec.Header
	if err := cc.ReadHeader(&h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("rpc server: read header error:", err)
		}
		return nil, err
	}
	return &h, nil
}

func (server *Server) readRequest(cc codec.Codec) (*request, error) {
	return nil, nil
}

func (server *Server) sendResponse(cc codec.Codec, h *codec.Header, body interface{}, locker *sync.Mutex) {
	locker.Lock()

	locker.Unlock()
}

func (server *Server) handleRequest(cc codec.Codec, req *request, locker *sync.Mutex, wg *sync.WaitGroup) {
	locker.Lock()

	locker.Unlock()
}
