package trpc

import (
	"bytes"
	"io"
	"net/http"
	"sync"
)

type HttpServer struct {
	c chan *RpcMessage
}

func NewHttpServer() (*HttpServer, error) {
	return &HttpServer{c: make(chan *RpcMessage, 100)}, nil
}
func (t *HttpServer) Handle(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	t.c <- &RpcMessage{Request: body, Writer: NewWriter(func(p []byte) (n int, err error) {
		defer wg.Done()
		return w.Write(p)
	})}
	wg.Wait()
}

func (t *HttpServer) Close() error {
	close(t.c)
	return nil
}
func (t *HttpServer) Get() <-chan *RpcMessage {
	return t.c
}

type HttpClient struct {
	conn *http.Client
	url  string
}

func NewHttpClient(url string) (*HttpClient, error) {
	conn := &http.Client{}
	t := HttpClient{conn: conn, url: url}
	return &t, nil
}

func (t *HttpClient) Request(data []byte) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, t.url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	resp, err := t.conn.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (t *HttpClient) Close() error {
	return nil
}
