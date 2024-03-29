package dwebsocket

import (
	"encoding/json"
	"fmt"
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
		fmt.Println(err)
	}

	data, err := json.Marshal(request.Data)
	if err != nil {
		fmt.Println(err)
		return
	}

	var resp []byte
	if request.Code >= 10000 && request.Code < 20000 {

	} else if request.Code >= 20000 && request.Code < 30000 {

	} else if request.Code >= 30000 && request.Code < 40000 {

	} else {
		// 采用 map 注册的方式
		value, ok := getHandlers(request.Code)
		if !ok {
			return
		}

		resp, err = value(client, data)
		if err != nil {
			return
		}
	}

	client.SendMsg(resp)
}
