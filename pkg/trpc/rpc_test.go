package trpc

import (
	"fmt"
	"testing"
)

func TestNewRpc(t *testing.T) {
	r, err := NewRpc(Option().
		//SetUrl("123").
		SetTransport(Http),
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(r)
}
