package dig

import (
	"db-go-game/domain/dao"
	"db-go-game/pkg/common/dgin"
	"db-go-game/services/api/internal/controller"
	"db-go-game/services/api/internal/service"
	"go.uber.org/dig"
)

var container = dig.New()

func init() {
	container.Provide(controller.NewUserController)
	container.Provide(service.NewUserService)
	container.Provide(dao.NewAuthRepository)
	container.Provide(dgin.NewGinServer)
}

func Invoke(i interface{}) error {
	return container.Invoke(i)
}
