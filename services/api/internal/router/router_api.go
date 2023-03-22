package router

import (
	"db-go-game/services/api/dig"
	"db-go-game/services/api/internal/controller"
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {
	publicGroup := engine.Group("open")
	registerPublicRoutes(publicGroup)

	privateGroup := engine.Group("api")
	registerPrivateRouter(privateGroup)
}

func registerPublicRoutes(group *gin.RouterGroup) {

}

func registerPrivateRouter(group *gin.RouterGroup) {
	var userController controller.IAuthController
	dig.Invoke(func(u controller.IAuthController) {
		userController = u
	})

	group.POST("login", userController.SignIn)
}
