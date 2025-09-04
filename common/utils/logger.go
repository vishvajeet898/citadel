package utils

import (
	"context"

	"github.com/Orange-Health/citadel/common/constants"
	"github.com/rs/zerolog/log"
)

func AddLog(ctx context.Context, level string, message string, params map[string]interface{}, err error) {
	if params == nil {
		params = map[string]interface{}{}
	}
	if traceID, ok := ctx.Value(constants.TraceIdKey).(string); ok {
		params["trace_id"] = traceID
	}
	switch level {
	case constants.DEBUG_LEVEL:
		logger := log.Debug().Fields(params)
		if err != nil {
			logger.Err(err).Msg(message)
		} else {
			logger.Msg(message)
		}
	case constants.INFO_LEVEL:
		logger := log.Info().Fields(params)
		if err != nil {
			logger.Err(err).Msg(message)
		} else {
			logger.Msg(message)
		}
	case constants.WARN_LEVEL:
		logger := log.Warn().Fields(params)
		if err != nil {
			logger.Err(err).Msg(message)
		} else {
			logger.Msg(message)
		}
	case constants.ERROR_LEVEL:
		logger := log.Error().Fields(params)
		if err != nil {
			logger.Err(err).Msg(message)
		} else {
			logger.Msg(message)
		}
	default:
		log.Warn().Msg("Unknown log level: " + level)
	}
}
