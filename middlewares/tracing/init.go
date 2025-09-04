package tracing

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := uuid.New().String()
		c.Set("trace_id", traceID)
		c.Next()
	}
}
