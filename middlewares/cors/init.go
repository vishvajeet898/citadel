package cors

import (
	"github.com/Orange-Health/citadel/conf"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Initialize() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: conf.GetConfig().GetStringSlice("cors_allow_origin"),
		AllowMethods: []string{"DELETE", "PATCH", "GET", "POST", "PATCH", "PUT", "OPTIONS"},
		AllowHeaders: []string{"Authorization", "content-type"},
	})
}
