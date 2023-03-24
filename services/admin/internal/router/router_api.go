package router

import (
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
	//var userController controller.IUserController
	//dig.Invoke(func(u controller.IUserController) {
	//	userController = u
	//})

}
