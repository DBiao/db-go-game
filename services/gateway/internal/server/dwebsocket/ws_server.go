package dwebsocket

import (
	"db-go-game/pkg/common/dgin"
	"db-go-game/pkg/conf"
	"db-go-game/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

const (
	// 最大的消息大小
	maxMessageSize = 8192
)

type WServer struct {
	port     int
	serverId int
	gin      *dgin.GinServer
}

func NewWServer(wsServer *conf.WsServer) *WServer {
	var (
		ws *WServer
	)
	ws = &WServer{
		port:     wsServer.Port,
		serverId: wsServer.ServerId,
		gin:      dgin.NewGinServer(),
	}
	ws.gin.Engine.Use(middleware.JwtAuth())
	ws.gin.Engine.GET("/", ws.wsHandler)
	return ws
}

func (ws *WServer) Run() {
	go func() {
		ws.gin.Run(ws.port)
	}()
}

func (ws *WServer) wsHandler(ctx *gin.Context) {
	conn, err := (&websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// 允许所有CORS跨域请求
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		http.NotFound(ctx.Writer, ctx.Request)
		return
	}

	//设置读取消息大小上线
	conn.SetReadLimit(maxMessageSize)

	// 验证token是否存在

	clientId := ctx.Param("clientId")
	cId, _ := strconv.Atoi(clientId)
	appId := ctx.Param("appId")
	aId, _ := strconv.Atoi(appId)
	client := newClient(conn, int64(cId), strconv.Itoa(aId))

	client.Close()

	go client.Run(conn)

	for {
		_, b, err := conn.ReadMessage()
		if err != nil {
			return
		}
		client.InChan <- b
	}
}
