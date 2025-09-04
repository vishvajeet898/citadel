package sentry

import (
	"context"

	"github.com/getsentry/sentry-go"

	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/utils"
)

type SentryLayer interface {
	LogError(ctx context.Context, message string, err error, additionalInformation map[string]interface{})
}

func getSentryHub() *sentry.Hub {
	hub := sentry.CurrentHub()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("service", "Citadel")
	})
	return hub
}

func (s *Sentry) LogError(ctx context.Context, message string, err error, additionalInformation map[string]interface{}) {
	if additionalInformation == nil {
		additionalInformation = make(map[string]interface{})
	}
	if ctx.Value("trace_id") != nil {
		additionalInformation["trace_id"] = ctx.Value("trace_id")
	}
	utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), additionalInformation, err)
	hub := getSentryHub()
	if hub != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			for key, value := range additionalInformation {
				scope.SetExtra(key, value)
			}
			if err == nil {
				hub.CaptureMessage(message)
			} else {
				hub.CaptureMessage(message + ": " + err.Error())
			}
		})
	}
}
