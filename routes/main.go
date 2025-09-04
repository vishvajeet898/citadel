package routes

import (
	"context"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/psql"
	"github.com/Orange-Health/citadel/adapters/sentry"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	corsMiddleware "github.com/Orange-Health/citadel/middlewares/cors"
	jwtMiddleware "github.com/Orange-Health/citadel/middlewares/jwt"
	loggerMiddleware "github.com/Orange-Health/citadel/middlewares/logger"
	tracingMiddleware "github.com/Orange-Health/citadel/middlewares/tracing"
)

type GinRouter struct {
	Gin *gin.Engine
}

func NewGinRouter() GinRouter {
	router := gin.New()
	ctx := context.Background()
	router.Use(gin.Recovery())
	psql.Initialize()
	cache.Initialize(ctx)
	router.Use(tracingMiddleware.Middleware())
	router.Use(jwtMiddleware.Middleware())
	router.Use(loggerMiddleware.Middleware())

	sentry.Initialize()
	router.Use(sentrygin.New(sentrygin.Options{Repanic: true}))

	router.Use(corsMiddleware.Initialize())

	router.Use(SetTraceIdInContext())

	initialiseRoutes(router)

	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	return GinRouter{
		Gin: router,
	}
}

func SetTraceIdInContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		if traceID, ok := c.Get("trace_id"); ok {
			c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), commonConstants.TraceIdKey, traceID))
		}
		c.Next()
	}
}
