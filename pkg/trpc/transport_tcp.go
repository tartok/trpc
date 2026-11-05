package trpc

import (
	"bufio"
	"net"
	"sync"
)

type TcpServer struct {
	listener net.Listener
	c        chan *RpcMessage
}

func NewTcpServer(address string) (*TcpServer, error) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	return &TcpServer{listener: l, c: make(chan *RpcMessage, 100)}, nil
}
func (t *TcpServer) Listen() error {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			return err
		}
		go func(conn net.Conn) {
			reader := bufio.NewReader(conn)
			w := NewWriter(func(p []byte) (n int, err error) {
				return conn.Write(append(p, 0))
			})
			for {
				data, err := reader.ReadBytes(0)
				if err != nil {
					return
				}
				t.c <- &RpcMessage{Request: data[:len(data)-1], Writer: w}
			}
		}(conn)
	}
}
func (t *TcpServer) Close() error {
	close(t.c)
	return t.listener.Close()
}
func (t *TcpServer) Get() <-chan *RpcMessage {
	return t.c
}

type TcpClient struct {
	sync.Mutex
	conn net.Conn
	c    chan []byte
}

func NewTcpClient(url string) (*TcpClient, error) {
	conn, err := net.Dial("tcp", url)
	if err != nil {
		return nil, err
	}
	t := TcpClient{conn: conn, c: make(chan []byte)}
	go func() {
		r := bufio.NewReader(conn)
		for {
			b, err := r.ReadBytes(0)
			if err != nil {
				return
			}
			t.c <- b[:len(b)-1]
		}
	}()
	return &t, nil
}
func (t *TcpClient) Write(p []byte) (n int, err error) {
	t.Lock()
	defer t.Unlock()
	return t.conn.Write(append(p, 0))
}

func (t *TcpClient) Receive() <-chan []byte {
	return t.c
}
func (t *TcpClient) Close() error {
	t.Lock()
	defer t.Unlock()
	close(t.c)
	return t.conn.Close()
}
