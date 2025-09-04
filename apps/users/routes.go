package users

import (
	"github.com/gin-gonic/gin"

	"github.com/Orange-Health/citadel/apps/users/controller"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	tokenAuthMiddleware "github.com/Orange-Health/citadel/middlewares/token_auth"
)

func RouteHandler(router *gin.RouterGroup) {
	userController := controller.InitUserController()

	router.GET("/", userController.GetAllUsers)
	router.GET("/:id", userController.GetUser)
	router.POST("/create-user", userController.CreateUser)
	router.PUT("/update-user/:id", userController.UpdateUser)
	router.DELETE("/delete-user/:id", userController.DeleteUser)
	router.GET("/get-user-by-email", tokenAuthMiddleware.Authenticate(commonConstants.PorteServiceName), userController.GetUserByEmail)
}
