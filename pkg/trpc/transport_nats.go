package trpc

import (
	"bytes"
	"io"
	"net/http"

	"github.com/nats-io/nats.go"
)

type NatsServer struct {
	c  chan *RpcMessage
	nc *nats.Conn
}

func NewNatsServer(nc *nats.Conn, subjectPrefix string, subjects ...string) (*NatsServer, error) {
	t := NatsServer{c: make(chan *RpcMessage, 100), nc: nc}
	if len(subjects) > 0 {
		for _, subject := range subjects {
			_, err := nc.Subscribe(subjectPrefix+"."+subject, func(msg *nats.Msg) {
				t.c <- &RpcMessage{Request: msg.Data, Writer: NewWriter(func(p []byte) (n int, err error) {
					err = msg.Respond(p)
					if err != nil {
						return 0, err
					}
					return len(p), nil
				})}
			})
			if err != nil {
				return nil, err
			}
		}
	} else {
		_, err := nc.Subscribe(subjectPrefix, func(msg *nats.Msg) {
			t.c <- &RpcMessage{Request: msg.Data, Writer: NewWriter(func(p []byte) (n int, err error) {
				err = msg.Respond(p)
				if err != nil {
					return 0, err
				}
				return len(p), nil
			})}
		})
		if err != nil {
			return nil, err
		}
	}
	return &t, nil
}

func (t *NatsServer) Close() error {
	close(t.c)
	t.nc.Close()
	return nil
}
func (t *NatsServer) Get() <-chan *RpcMessage {
	return t.c
}

type NatsClient struct {
	conn *http.Client
	url  string
}

func NewNatsClient(url string) (*NatsClient, error) {
	conn := &http.Client{}
	t := NatsClient{conn: conn, url: url}
	return &t, nil
}

func (t *NatsClient) Request(data []byte) ([]byte, error) {
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

func (t *NatsClient) Close() error {
	return nil
}
