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
			for {
				data, index, err := ReadData(reader)
				if err != nil {
					return
				}
				w := NewWriter(func(p []byte) (n int, err error) {
					err = WriteData(conn, index, p)
					return len(p), err
				})

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

type IndexChan struct {
	C    chan []byte
	Busy bool
}
type TcpClient struct {
	sync.Mutex
	conn net.Conn
	//c    chan []byte
	c []*IndexChan
}

func NewTcpClient(url string) (*TcpClient, error) {
	conn, err := net.Dial("tcp", url)
	if err != nil {
		return nil, err
	}
	t := TcpClient{conn: conn}
	go func() {
		r := bufio.NewReader(conn)
		for {
			b, index, err := ReadData(r)
			if err != nil {
				_ = t.Close()
				return
			}
			//fmt.Println("read:", index)
			t.Lock()
			c := t.getC(index)
			if c != nil {
				c <- b
			}
			t.Unlock()
		}
	}()
	return &t, nil
}
func (t *TcpClient) getC(index int) chan []byte {
	if index >= len(t.c) {
		return nil
	}
	i := t.c[index]
	i.Busy = false
	return i.C
}
func (t *TcpClient) newC() (int, chan []byte) {
	t.Lock()
	defer t.Unlock()
	for i, c := range t.c {
		if !c.Busy {
			c.Busy = true
			return i, c.C
		}
	}
	i := len(t.c)
	c := make(chan []byte)
	t.c = append(t.c, &IndexChan{C: c, Busy: true})
	return i, c
}
func (t *TcpClient) write(index int, p []byte) (n int, err error) {
	t.Lock()
	defer t.Unlock()
	err = WriteData(t.conn, index, p)
	return len(p), err
}

func (t *TcpClient) Request(data []byte) (resp []byte, err error) {
	index, c := t.newC()
	//fmt.Println("newC", index)
	_, err = t.write(index, data)
	if err != nil {
		return
	}
	r, ok := <-c
	if !ok {
		return nil, errors.New("transport closed")
	}
	return r, nil
}

func (t *TcpClient) Close() error {
	t.Lock()
	defer t.Unlock()
	for _, indexChan := range t.c {
		close(indexChan.C)
	}
	t.c = t.c[:0]
	return t.conn.Close()
}
