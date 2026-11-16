package trpc

import (
	"encoding/json"
	"time"

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
	nc            *nats.Conn
	subjectPrefix string
	timeout       time.Duration
}

func NewNatsClient(nc *nats.Conn, subjectPrefix string, timeout time.Duration) (*NatsClient, error) {
	t := NatsClient{nc: nc, subjectPrefix: subjectPrefix, timeout: timeout}
	return &t, nil
}

func (t *NatsClient) Request(data []byte) ([]byte, error) {
	var m struct {
		Method string `json:"method"`
	}
	err := json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	msg, err := t.nc.Request(t.subjectPrefix+"."+m.Method, data, t.timeout)
	if err != nil {
		return nil, err
	}
	return msg.Data, nil
	//req, err := http.NewRequest(http.MethodPost, t.url, bytes.NewReader(data))
	//if err != nil {
	//	return nil, err
	//}
	//resp, err := t.conn.Do(req)
	//if err != nil {
	//	return nil, err
	//}
	//defer resp.Body.Close()
	//body, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	return nil, err
	//}
	//return body, nil
}

func (t *NatsClient) Close() error {
	t.nc.Close()
	return nil
}
