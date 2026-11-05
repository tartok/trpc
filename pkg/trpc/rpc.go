package trpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Rpc struct {
	opt Options
}

func NewRpc(o ...*OptionsBuilder) (*Rpc, error) {
	var r Rpc
	n, err := NewOptions(o...)
	if err != nil {
		return nil, err
	}
	r.opt = *n
	return &r, nil
}

func (r *Rpc) HttpRequest(id Id, method string, params any, res any, url, auth string) (*Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(NewRequest(id, method, params).Bytes()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("authorization", auth)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	var result Response
	if res != nil {
		result.Result = res
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil

}
