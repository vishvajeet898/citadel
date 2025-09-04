package server

import (
	"github.com/Orange-Health/citadel/conf"
	"github.com/Orange-Health/citadel/routes"
)

func Start() {
	config := conf.GetConfig()

	router := routes.NewGinRouter()
	router.Gin.Run(config.GetString("server.host") + ":" + config.GetString("server.port"))
}
