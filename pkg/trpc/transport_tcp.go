package trpc

import (
	"bufio"
	"errors"
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
			var index int
			w := NewWriter(func(p []byte) (n int, err error) {
				err = WriteData(conn, index, p)
				return len(p), err
			})
			var data []byte
			var err error
			for {
				data, index, err = ReadData(reader)
				if err != nil {
					return
				}
				t.c <- &RpcMessage{Request: data, Writer: w}
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
			b, _, err := ReadData(r)
			if err != nil {
				return
			}
			t.Lock()
			if t.c != nil {
				t.c <- b
			}
			t.Unlock()
		}
	}()
	return &t, nil
}
func (t *TcpClient) write(p []byte) (n int, err error) {
	t.Lock()
	defer t.Unlock()
	err = WriteData(t.conn, 0, p)
	return len(p), err
}

func (t *TcpClient) Request(data []byte) (resp []byte, err error) {
	_, err = t.write(data)
	if err != nil {
		return
	}
	r, ok := <-t.c
	if !ok {
		return nil, errors.New("transport closed")
	}
	return r, nil
}

func (t *TcpClient) Close() error {
	t.Lock()
	defer t.Unlock()
	close(t.c)
	t.c = nil
	return t.conn.Close()
}
