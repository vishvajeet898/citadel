package reportRebrandingClient

import (
	"github.com/Orange-Health/citadel/conf"
)

var (
	Config                         = conf.GetConfig()
	ReportRebrandingBaseUrl string = Config.GetString("report_rebranding.base_url")
	ReportRebrandingToken   string = Config.GetString("report_rebranding.token")
)

const (
	MEDIA_RESIZE           = "MEDIA_RESIZE"
	ATTACH_COBRANDED_IMAGE = "ATTACH_COBRANDED_IMAGE"
)

var URL_MAP = map[string]string{
	MEDIA_RESIZE:           "/media-resize/resize",
	ATTACH_COBRANDED_IMAGE: "/pdf-gen/attach-cobranded-image",
}
