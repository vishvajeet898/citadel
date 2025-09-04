package slackClient

import (
	"github.com/Orange-Health/citadel/conf"
)

var (
	Config       = conf.GetConfig()
	SlackBaseUrl = Config.GetString("slack.base_url")
	SlackToken   = Config.GetString("slack.token")
)

const (
	CHAT_POST_MESSAGE = "/api/chat.postMessage"
)

var URL_MAP = map[string]string{
	"post_message": CHAT_POST_MESSAGE,
}
