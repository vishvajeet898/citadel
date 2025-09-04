package utils

import (
	"github.com/Orange-Health/citadel/common/constants"
	pubsub "github.com/Orange-Health/pubsublib/provider/aws"
)

func GetNewPubSubAdapter(accessKeyId, secretAccessKey string) (*pubsub.AWSPubSubAdapter, error) {
	awsAdapter, err := pubsub.NewAWSPubSubAdapter(
		constants.Region,
		accessKeyId,
		secretAccessKey,
		"",
		constants.PubSubRedisAddr,
		constants.PubSubRedisPass,
		constants.PubSubRedisDb,
		constants.PubSubRedisPoolSize,
		constants.PubSubRedisMinIdleConn,
	)
	return awsAdapter, err
}
