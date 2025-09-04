package sentry

import (
	"time"

	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/conf"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog/log"
)

type Sentry struct{}

func Initialize() {
	config := conf.GetConfig()
	sentryDsn := config.GetString("sentry.dsn")
	if sentryDsn != "" {
		log.Info().Msg("setting up sentry")
		// To initialize Sentry's handler, you need to initialize Sentry itself beforehand
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:              sentryDsn,
			AttachStacktrace: true,
		}); err != nil {
			log.Error().Err(err).Msg(constants.SENTRY_INITIALIZATION_FAILED)
		}
	}
	defer sentry.Flush(time.Second * 2)
	defer sentry.Recover()
}

func GetSentryInstance() *Sentry {
	return &Sentry{}
}

func InitializeSentry() SentryLayer {
	return GetSentryInstance()
}
