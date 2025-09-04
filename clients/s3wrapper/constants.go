package s3wrapperClient

import (
	"github.com/Orange-Health/citadel/conf"
)

var (
	Config           = conf.GetConfig()
	S3wrapperBaseUrl = Config.GetString("s3wrapper.base_url")
	S3wrapperApiKey  = Config.GetString("s3wrapper.api_key")
)

const (
	ORDER_FILE_PUBLIC_URL = "ORDER_FILE_PUBLIC_URL"
)

var URL_MAP = map[string]string{
	ORDER_FILE_PUBLIC_URL: "/api/media-url/",
}
