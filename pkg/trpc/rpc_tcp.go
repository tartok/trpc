package trpc

import (
	"bufio"
	"encoding/json"
	"net"
)

func (r *RpcClient) TcpRequest(method string, params any, res any, customOptions ...Lister[ClientOptions]) (*Response, error) {
	options := r.opt
	err := AddOptions(&options, customOptions...)

	conn, err := net.Dial("tcp", options.Url)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = conn.Close()
	}()

	request := NewRequest(nil, method, params)
	payload := request.Bytes()
	payload = append(payload, 0)
	_, err = conn.Write(payload)
	if err != nil {
		return nil, err
	}

	buff, err := bufio.NewReader(conn).ReadBytes(0)
	if err != nil {
		return nil, err
	}
	var result Response
	if res != nil {
		result.Result = res
	}
	err = json.Unmarshal(buff[:len(buff)-1], &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
