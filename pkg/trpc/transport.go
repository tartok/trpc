package trpc

import (
	"io"
	"sync"
)

type RpcMessage struct {
	Request []byte
	Writer  io.Writer
}

type TransportServer interface {
	Get() <-chan *RpcMessage
}
type TransportClient interface {
	Request(data []byte) (resp []byte, err error)
}

type Writer struct {
	sync.Mutex
	wf func(p []byte) (n int, err error)
}

func (w *Writer) Write(p []byte) (n int, err error) {
	w.Lock()
	defer w.Unlock()
	return w.wf(p)
}

func NewWriter(wf func(p []byte) (n int, err error)) *Writer {
	return &Writer{wf: wf}
}
