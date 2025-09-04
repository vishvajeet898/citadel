package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	commonStructures "github.com/Orange-Health/citadel/common/structures"
)

func HandleORMError(err error) *commonStructures.CommonError {
	if err == gorm.ErrRecordNotFound {
		return &commonStructures.CommonError{
			Message:    err.Error(),
			StatusCode: http.StatusNotFound,
		}
	} else {
		return &commonStructures.CommonError{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}
}

func HandleErrorResponse(c *gin.Context, statusCode int, message string) {
	apiResponse := commonStructures.CommonAPIResponse{
		Error: message,
	}
	c.JSON(statusCode, apiResponse)
}
