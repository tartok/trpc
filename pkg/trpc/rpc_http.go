package trpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (r *RpcClient) HttpRequest(method string, params any, res any, customOptions ...Lister[ClientOptions]) (*Response, error) {
	options := r.opt
	err := AddOptions(&options, customOptions...)

	req, err := http.NewRequest(http.MethodPost, options.Url, bytes.NewReader(NewRequest(nil, method, params).Bytes()))
	if err != nil {
		return nil, err
	}
	if options.Auth != "" {
		req.Header.Add("authorization", options.Auth)
	}
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
