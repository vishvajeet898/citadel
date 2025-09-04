package health

import (
	"github.com/gin-gonic/gin"
)

func RouteHandler(router *gin.RouterGroup) {
	router.GET("/", HealthCheckController)
}
