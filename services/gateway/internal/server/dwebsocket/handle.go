package dwebsocket

import (
	"encoding/json"
	"sync"
)

// WsRequest 通用请求数据格式
type WsRequest struct {
	Code uint16      `json:"seq"`            // 消息的唯一Id
	Data interface{} `json:"data,omitempty"` // 数据 json
}

// WsResponse 通用返回数据格式
type WsResponse struct {
	Code uint32      `json:"seq"`            // 消息的唯一Id
	Data interface{} `json:"data,omitempty"` // 数据 json
}

type HandleFunc func(client *Client, message []byte) ([]byte, error)

var (
	handlers        = make(map[uint16]HandleFunc)
	handlersRWMutex sync.RWMutex
)

// Register 注册
func Register(key uint16, value HandleFunc) {
	handlersRWMutex.Lock()
	defer handlersRWMutex.Unlock()
	handlers[key] = value

	return
}

func getHandlers(code uint16) (value HandleFunc, ok bool) {
	handlersRWMutex.RLock()
	defer handlersRWMutex.RUnlock()

	value, ok = handlers[code]

	return
}

// Handle 处理数据
func Handle(client *Client, message []byte) {
	//defer utils.PrintPanic()

	request := &WsRequest{}
	err := json.Unmarshal(message, request)
	if err != nil {
	}

	data, err := json.Marshal(request.Data)
	if err != nil {
		return
	}

	switch request.Code {

	}

	// 采用 map 注册的方式
	value, ok := getHandlers(request.Code)
	if !ok {
		return
	}

	resp, err := value(client, data)
	if err != nil {
		return
	}

	client.SendMsg(resp)
}

func forward(code int32) {
	if code >= 10000 && code < 20000 {

	} else if code >= 20000 && code < 30000 {

	} else if code >= 30000 && code < 40000 {

	}
}
