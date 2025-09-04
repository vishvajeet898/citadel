package sqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/rs/zerolog/log"

	"github.com/Orange-Health/citadel/common/constants"
)

func getAWSSession() *session.Session {
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(constants.Region),
			Credentials: credentials.NewStaticCredentials(
				constants.SQSAccessKeyID,
				constants.SQSSecretAccessKey,
				"", // a token will be created when the session it's used.
			),
		})
	if err != nil {
		panic(err)
	}
	return sess
}

func GetSqsClient() *sqs.SQS {
	sess := getAWSSession()
	sqsSvc := sqs.New(sess)
	return sqsSvc
}

func InitAndGetSQSClient() *sqs.SQS {
	sqsSvc := GetSqsClient()
	_, err := sqsSvc.CreateQueue(&sqs.CreateQueueInput{
		Attributes: map[string]*string{
			"ReceiveMessageWaitTimeSeconds": aws.String("20"),
			"SqsManagedSseEnabled":          aws.String("true"),
		},
		QueueName: &constants.WorkerDefaultQueue,
	})
	if err != nil {
		log.Debug().Err(err).Msg(constants.ERROR_CREATING_QUEUE)
	}
	return sqsSvc
}
