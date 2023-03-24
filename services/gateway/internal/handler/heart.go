package handler

import "db-go-game/services/gateway/internal/server/dwebsocket"

func init() {
	dwebsocket.Register(1001, Heart)
}

func Heart(client *dwebsocket.Client, message []byte) ([]byte, error) {
	return nil, nil
}
