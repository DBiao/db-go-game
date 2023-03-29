package dwebsocket

import (
	"github.com/gorilla/websocket"
	"sync/atomic"
	"time"
)

type Client struct {
	//rwLock    sync.RWMutex
	conn      *websocket.Conn
	uid       int64  // 用户ID
	appId     string // appID
	onlineTs  int64  // 上线时间戳（毫秒）
	OutChan   chan []byte
	InChan    chan []byte
	QuitChan  chan struct{}
	CloseFlag int32
}

func newClient(conn *websocket.Conn, uid int64, appId string) *Client {
	return &Client{
		conn:     conn,
		uid:      uid,
		appId:    appId,
		onlineTs: time.Now().UnixNano() / 1e6,
		OutChan:  make(chan []byte, 1),
		InChan:   make(chan []byte, 1),
		QuitChan: make(chan struct{}, 1),
	}
}

func (c *Client) GetUid() int64 {
	return c.uid
}

func (c *Client) Run(ws *websocket.Conn) {
	//defer utils.PrintPanic()

	for {
		select {
		case in, ok := <-c.InChan:
			if !ok {
				return
			}
			Handle(c, in)
		case out, ok := <-c.OutChan:
			if !ok {
				return
			}
			err := ws.WriteMessage(websocket.TextMessage, out)
			if err != nil {
				return
			}
		case <-c.QuitChan:
			return
		}
	}
}

// Close 关闭监听
func (c *Client) Close() {
	c.conn.SetCloseHandler(func(code int, text string) error {
		c.stop()
		c.conn.Close()
		ClientMap.Delete(c.uid)
		return nil
	})
}

func (c *Client) stop() {
	if atomic.CompareAndSwapInt32(&c.CloseFlag, 0, 1) {
		select {
		case c.QuitChan <- struct{}{}:
		default:
		}
	}
}

// SendMsg 发送数据
func (c *Client) SendMsg(msg []byte) {
	if c == nil {
		return
	}

	c.OutChan <- msg
}
