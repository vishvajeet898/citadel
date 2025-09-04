package health

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Orange-Health/citadel/common/structures"
)

//	HealthCheck godoc
//
// @Summary	Health Check
// @Produce	json
// @Success	200	{object}	structures.CommonAPIResponse	"Common API Response"
// @Router		/ping [get]
func HealthCheckController(c *gin.Context) {
	c.JSON(http.StatusOK, structures.CommonAPIResponse{
		Message: "successful",
	})
}
