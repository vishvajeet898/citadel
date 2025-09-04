package tokenAuth

import (
	"net/http"

	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/conf"
	"github.com/gin-gonic/gin"
)

func Authenticate(service string) gin.HandlerFunc {
	return func(c *gin.Context) {
		validateToken(c, service)
		if !c.IsAborted() {
			c.Next()
		}
	}
}

func MutipleAuthenticate(services []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, service := range services {
			if multipleValidateToken(c, service) {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		c.Abort()
	}
}

func validateToken(c *gin.Context, service string) {
	token := c.Request.Header.Get("api-key")
	serviceKey := c.Request.Header.Get("service")
	config := conf.GetConfig()
	newService := service
	if service == constants.HealthServiceName {
		newService = service + "_api"
	}
	apiToken := config.GetString(newService + ".incoming_api_key")

	if token == apiToken && serviceKey == service {
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{
		"message": "Unauthorized",
	})
	c.Abort()
}

func multipleValidateToken(c *gin.Context, service string) bool {
	token := c.Request.Header.Get("api-key")
	serviceKey := c.Request.Header.Get("service")
	config := conf.GetConfig()
	apiToken := config.GetString(service + ".incoming_api_key")

	if token == apiToken && serviceKey == service {
		return true
	}

	return false
}
