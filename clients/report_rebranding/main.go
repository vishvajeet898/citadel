package reportRebrandingClient

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/clients"
	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/structures"
	"github.com/Orange-Health/citadel/common/utils"
)

type ReportRebrandingClient struct {
	ApiClient *clients.ApiClient
	Headers   map[string]string
	Cache     cache.CacheLayer
}

func NewClient() *ReportRebrandingClient {
	apiClient := clients.NewClient()
	apiClient.BaseURL = ReportRebrandingBaseUrl
	return &ReportRebrandingClient{
		ApiClient: apiClient,
		Headers: map[string]string{
			"service": constants.CitadelServiceName,
			"token":   ReportRebrandingToken,
		},
		Cache: cache.InitializeCache(),
	}
}

func InitializeReportRebrandingClient() ReportRebrandingClientInterface {
	return NewClient()
}

type ReportRebrandingClientInterface interface {
	ResizeMedia(ctx context.Context, mediaResizeRequest structures.MediaResizeRequest) (
		structures.MediaResizeResponse, error)
	AttachCobrandedImage(ctx context.Context, cobrandingRequest structures.CobrandingRequest) (
		structures.CobrandingResponse, error)
}

func (reportRebrandingClient *ReportRebrandingClient) ResizeMedia(ctx context.Context,
	mediaResizeRequest structures.MediaResizeRequest) (structures.MediaResizeResponse, error) {
	mediaResizeResponse := structures.MediaResizeResponse{}

	var response interface{}

	headers := reportRebrandingClient.Headers
	headers["Content-Type"] = constants.ContentTypeJson

	err := reportRebrandingClient.ApiClient.Post(ctx, &response, URL_MAP[MEDIA_RESIZE], nil, mediaResizeRequest, headers, 3, time.Millisecond*500)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_RESIZING_MEDIA, nil, err)
		return mediaResizeResponse, err
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_RESIZING_MEDIA, nil, err)
		return mediaResizeResponse, err
	}

	err = json.Unmarshal(responseBytes, &mediaResizeResponse)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_RESIZING_MEDIA, nil, err)
		return mediaResizeResponse, err
	}

	return mediaResizeResponse, nil
}

func (reportRebrandingClient *ReportRebrandingClient) AttachCobrandedImage(ctx context.Context,
	cobrandingRequest structures.CobrandingRequest) (structures.CobrandingResponse, error) {

	var response interface{}
	headers := reportRebrandingClient.Headers
	headers["Content-Type"] = constants.ContentTypeJson

	err := reportRebrandingClient.ApiClient.Post(ctx, &response, URL_MAP[ATTACH_COBRANDED_IMAGE], nil, cobrandingRequest, headers, 1, 0)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_ATTACHING_COBRANDED_IMAGE, nil, err)
		return structures.CobrandingResponse{}, err
	}

	cobrandingResponse := structures.CobrandingResponse{}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_ATTACHING_COBRANDED_IMAGE, nil, err)
		return cobrandingResponse, err
	}

	err = json.Unmarshal(responseBytes, &cobrandingResponse)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_ATTACHING_COBRANDED_IMAGE, nil, err)
		return cobrandingResponse, err
	}

	return cobrandingResponse, nil
}
