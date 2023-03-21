package dig

import (
	"db-go-game/services/logic/internal/server"
	"db-go-game/services/logic/internal/service"
	"go.uber.org/dig"
)

var container = dig.New()

func init() {
	container.Provide(service.NewAuthService)
	container.Provide(server.NewServer)
	container.Provide(server.NewGrpcServer)
}

func Invoke(i interface{}) error {
	return container.Invoke(i)
}
