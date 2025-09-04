package scheduler

import (
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/Orange-Health/citadel/adapters/sentry"
	"github.com/Orange-Health/citadel/adapters/sqs"
	"github.com/Orange-Health/citadel/common/constants"
)

func StartServer(clientMode bool) (*machinery.Server, error) {
	sqsClient := sqs.InitAndGetSQSClient()
	redisMaxActive := 1
	redisMaxIdle := 0
	if !clientMode {
		redisMaxActive = constants.WorkerRedisPoolSize
		redisMaxIdle = constants.WorkerMaxIdle
	}
	var visibilityTimeout = 43200
	var cnf = &config.Config{
		Broker:        constants.WorkerBroker,
		DefaultQueue:  constants.WorkerDefaultQueue,
		ResultBackend: constants.WorkerResultBackend,
		Redis: &config.RedisConfig{
			MaxActive: redisMaxActive,
			MaxIdle:   redisMaxIdle,
		},
		SQS: &config.SQSConfig{
			Client:            sqsClient,
			WaitTimeSeconds:   20,
			VisibilityTimeout: &visibilityTimeout,
		},
	}

	if constants.Environment == "local" {
		cnf.AMQP = &config.AMQPConfig{
			Exchange:     "machinery_exchange",
			ExchangeType: "direct",
			BindingKey:   "machinery_task",
		}
	}

	server, err := machinery.NewServer(cnf)
	if err != nil {
		return nil, err
	}

	return server, err
}

func Start() error {
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// initializing sentry
	sentry.Initialize()

	taskServer, err := StartServer(false)
	if err != nil {
		log.ERROR.Fatalf("Error starting task server: %v\n", err.Error())
		return err
	}

	HandlePeriodicTasks(taskServer)

	select {}
}
