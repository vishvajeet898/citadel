package attuneClient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
	"github.com/Orange-Health/citadel/clients"
	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/structures"
	"github.com/Orange-Health/citadel/common/utils"
)

type AttuneClient struct {
	ApiClient *clients.ApiClient
	Headers   map[string]string
	Cache     cache.CacheLayer
	Sentry    sentry.SentryLayer
}

func NewClient() *AttuneClient {
	apiClient := clients.NewClient()
	apiClient.BaseURL = constants.AttuneAPIBaseUrl
	return &AttuneClient{
		ApiClient: apiClient,
		Headers:   map[string]string{},
		Cache:     cache.InitializeCache(),
		Sentry:    sentry.InitializeSentry(),
	}
}

func InitializeAttuneClient() AttuneClientInterface {
	return NewClient()
}

type AttuneClientInterface interface {
	GetPatientVisitDetailsbyVisitNo(ctx context.Context, visitId, reportPdfFormat string, labId uint) (
		structures.AttuneOrderResponse, *structures.CommonError)
	InsertTestDataToAttune(ctx context.Context, attuneOrder structures.AttuneOrderResponse) *structures.CommonError
	SyncDataToAttune(ctx context.Context, labId uint,
		payload structures.AttuneSyncDataToLisRequest) *structures.CommonError
	SyncDataToAttuneAfterSync(ctx context.Context, labId uint,
		payload structures.AttuneSyncDataToLisAfterSyncRequest) *structures.CommonError
}

func (attuneClient *AttuneClient) login(ctx context.Context) string {
	var resp interface{}

	headers := attuneClient.Headers
	headers["Content-Type"] = "text/plain"

	err := attuneClient.ApiClient.Get(ctx, &resp, ATTUNE_URLS[LOGIN], nil, attuneLoginPayload, headers, LoginRetries,
		PatientResultDetailsDelay)
	if err != nil {
		attuneClient.Sentry.LogError(ctx, constants.ERROR_WHILE_LOGGING_INTO_ATTUNE, err, nil)
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_LOGGING_INTO_ATTUNE, nil, err)
		return ""
	}

	response := structures.AttuneAuthResponse{}
	responseBytes, err := json.Marshal(resp)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_LOGGING_INTO_ATTUNE, nil, err)
		return ""
	}

	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_LOGGING_INTO_ATTUNE, nil, err)
		return ""
	}

	AttuneAccessToken = response.AccessToken
	AttuneTokenExpiry = response.ExpiresIn
	AttuneTokenCreatedOn = time.Now().Unix()

	return response.AccessToken
}

func (attuneCLient *AttuneClient) getToken(ctx context.Context) string {
	timeNow := time.Now().Unix()
	token := AttuneAccessToken

	if token == "" || ((timeNow - AttuneTokenCreatedOn) > 20000) {
		// refresh token if 20k seconds left
		token = attuneCLient.login(ctx)
	}
	return token
}

func fetchAttuneOrderResponseAfterFetchingVisitDetails(ctx context.Context, response interface{}) (
	structures.AttuneOrderResponse, *structures.CommonError) {
	responseBytes, err := json.Marshal(response)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_GETTING_ORDER_DATA_FROM_ATTUNE, nil, err)
		return structures.AttuneOrderResponse{}, &structures.CommonError{
			Message:    constants.ERROR_WHILE_GETTING_ORDER_DATA_FROM_ATTUNE,
			StatusCode: http.StatusInternalServerError,
		}
	}

	var attuneSyncResponse structures.AttuneOrderResponse
	err = json.Unmarshal([]byte(responseBytes), &attuneSyncResponse)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_GETTING_ORDER_DATA_FROM_ATTUNE, nil, err)
		return structures.AttuneOrderResponse{}, &structures.CommonError{
			Message:    constants.ERROR_WHILE_GETTING_ORDER_DATA_FROM_ATTUNE,
			StatusCode: http.StatusInternalServerError,
		}
	}

	if attuneSyncResponse.OrderId == "" {
		responseBytes, err := json.Marshal(response)
		if err != nil {
			utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_GETTING_ORDER_DATA_FROM_ATTUNE, nil, err)
			return structures.AttuneOrderResponse{}, &structures.CommonError{
				Message:    constants.ERROR_WHILE_GETTING_ORDER_DATA_FROM_ATTUNE,
				StatusCode: http.StatusInternalServerError,
			}
		}

		var attuneGetSyncResponse structures.AttuneGetSyncDataResponse
		err = json.Unmarshal([]byte(responseBytes), &attuneGetSyncResponse)
		if err != nil {
			utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_GETTING_ORDER_DATA_FROM_ATTUNE, nil, err)
			return structures.AttuneOrderResponse{}, &structures.CommonError{
				Message:    constants.ERROR_WHILE_GETTING_ORDER_DATA_FROM_ATTUNE,
				StatusCode: http.StatusInternalServerError,
			}
		}

		return attuneGetSyncResponse.Response, nil
	}

	return attuneSyncResponse, nil
}

func (attuneClient *AttuneClient) GetPatientVisitDetailsbyVisitNo(ctx context.Context,
	visitId, reportPdfFormat string, labId uint) (
	structures.AttuneOrderResponse, *structures.CommonError) {
	loggingAttributes := map[string]interface{}{}
	token := attuneClient.getToken(ctx)
	if token == "" {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_GETTING_ORDER_DATA_FROM_ATTUNE, nil, nil)
		return structures.AttuneOrderResponse{}, &structures.CommonError{
			Message:    constants.ERROR_WHILE_GETTING_ORDER_DATA_FROM_ATTUNE,
			StatusCode: http.StatusInternalServerError,
		}
	}

	if !utils.SliceContainsString(constants.AttuneReportdfTypes, reportPdfFormat) {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_GETTING_ORDER_DATA_FROM_ATTUNE, nil, nil)
		return structures.AttuneOrderResponse{}, &structures.CommonError{
			Message:    constants.ERROR_INVALID_REPORT_PDF_FORMAT,
			StatusCode: http.StatusBadRequest,
		}
	}

	orgCode := utils.GetAttuneOrgCodeByLabId(labId)
	utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), map[string]interface{}{
		"visit_id":          visitId,
		"report_pdf_format": reportPdfFormat,
		"city_code":         labId,
		"org_code":          orgCode,
	}, nil)
	if orgCode == "" {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_GETTING_ORDER_DATA_FROM_ATTUNE, nil, nil)
		return structures.AttuneOrderResponse{}, &structures.CommonError{
			Message:    constants.ERROR_GETTING_ATTUNE_ORG_CODE,
			StatusCode: http.StatusInternalServerError,
		}
	}
	queryParams := map[string]interface{}{
		"OrgCode":         orgCode,
		"VisitNumber":     visitId,
		"ReportPDFFormat": reportPdfFormat,
	}

	var response interface{}
	headers := attuneClient.Headers
	headers["Authorization"] = fmt.Sprintf("Bearer %s", token)
	headers["Content-Type"] = constants.ContentTypeJson

	err := attuneClient.ApiClient.Get(ctx, &response, ATTUNE_URLS[GET_PATIENT_RESULT_BY_VISIT], queryParams, nil, headers,
		PatientResultDetailsRetries, PatientResultDetailsDelay)
	loggingAttributes["response"] = response
	utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), loggingAttributes, nil)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_GETTING_ORDER_DATA_FROM_ATTUNE, loggingAttributes,
			err)
		return structures.AttuneOrderResponse{}, &structures.CommonError{
			Message:    constants.ERROR_WHILE_GETTING_ORDER_DATA_FROM_ATTUNE,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return fetchAttuneOrderResponseAfterFetchingVisitDetails(ctx, response)
}

func (attuneClient *AttuneClient) InsertTestDataToAttune(ctx context.Context,
	attuneOrder structures.AttuneOrderResponse) *structures.CommonError {

	loggingAttributes := map[string]interface{}{}
	token := attuneClient.getToken(ctx)
	if token == "" {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_SENDING_ATTUNE_WEBHOOK, nil, nil)
		return &structures.CommonError{
			Message:    constants.ERROR_WHILE_SENDING_ATTUNE_WEBHOOK,
			StatusCode: http.StatusInternalServerError,
		}
	}

	headers := attuneClient.Headers
	headers["Authorization"] = fmt.Sprintf("Bearer %s", token)
	headers["Content-Type"] = constants.ContentTypeJson

	utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), map[string]interface{}{
		"attune_order": attuneOrder,
	}, nil)

	var response interface{}
	err := attuneClient.ApiClient.Post(ctx, &response, ATTUNE_URLS[INSERT_DOCTOR_DATA_TO_ATTUNE], nil, attuneOrder, headers,
		SyncDataToAttuneRetries, SyncDataToAttuneDelay)
	loggingAttributes["response"] = response
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_SENDING_ATTUNE_WEBHOOK, loggingAttributes, err)
		return &structures.CommonError{
			Message:    constants.ERROR_WHILE_SENDING_ATTUNE_WEBHOOK,
			StatusCode: http.StatusInternalServerError,
		}
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_SENDING_ATTUNE_WEBHOOK, loggingAttributes, err)
		return &structures.CommonError{
			Message:    constants.ERROR_WHILE_SENDING_ATTUNE_WEBHOOK,
			StatusCode: http.StatusInternalServerError,
		}
	}

	attuneSyncResponse := structures.AttuneSyncResponse{}
	err = json.Unmarshal(responseBytes, &attuneSyncResponse)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_SENDING_ATTUNE_WEBHOOK, loggingAttributes, err)
		return &structures.CommonError{
			Message:    constants.ERROR_WHILE_SENDING_ATTUNE_WEBHOOK,
			StatusCode: http.StatusInternalServerError,
		}
	}

	if attuneSyncResponse.Code == "Error" {
		loggingAttributes["code"] = attuneSyncResponse.Code
		loggingAttributes["message"] = attuneSyncResponse.Message
		loggingAttributes["visit_id"] = attuneOrder.OrderId
		utils.AddLog(ctx, constants.ERROR_LEVEL, attuneSyncResponse.Message, loggingAttributes, nil)
		attuneClient.Sentry.LogError(ctx, attuneSyncResponse.Message, nil, loggingAttributes)
		return &structures.CommonError{
			Message:    attuneSyncResponse.Message,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}

func (attuneClient *AttuneClient) SyncDataToAttune(ctx context.Context, labId uint,
	payload structures.AttuneSyncDataToLisRequest) *structures.CommonError {

	token := attuneClient.getToken(ctx)
	if token == "" {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_SENDING_ATTUNE_WEBHOOK, nil, nil)
		return &structures.CommonError{
			Message:    constants.ERROR_WHILE_SENDING_ATTUNE_WEBHOOK,
			StatusCode: http.StatusInternalServerError,
		}
	}

	orgCode := utils.GetAttuneOrgCodeByLabId(labId)
	if orgCode == "" {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_SENDING_ATTUNE_WEBHOOK, nil, nil)
		return &structures.CommonError{
			Message:    constants.ERROR_GETTING_ATTUNE_ORG_CODE,
			StatusCode: http.StatusInternalServerError,
		}
	}

	headers := attuneClient.Headers
	headers["Authorization"] = fmt.Sprintf("Bearer %s", token)
	headers["IsWrapReq"] = "Y"
	headers["OrgCode"] = orgCode
	headers["Content-Type"] = constants.ContentTypeJson

	utils.AddLog(ctx, constants.INFO_LEVEL, utils.GetCurrentFunctionName(), map[string]interface{}{
		"payload": payload,
	}, nil)

	var response interface{}
	err := attuneClient.ApiClient.Post(ctx, &response, ATTUNE_URLS[SYNC_DATA_TO_ATTUNE], nil, payload, headers,
		SyncDataToAttuneRetries, SyncDataToAttuneDelay)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_SENDING_ATTUNE_WEBHOOK, nil, err)
		return &structures.CommonError{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_SENDING_ATTUNE_WEBHOOK, nil, err)
		return &structures.CommonError{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	attuneBookingResponse := structures.AttuneBookingResponse{}
	err = json.Unmarshal([]byte(responseBytes), &attuneBookingResponse)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_SENDING_ATTUNE_WEBHOOK, nil, err)
		return &structures.CommonError{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	if strings.ToLower(attuneBookingResponse.Response.Code) != "success" {
		utils.AddLog(ctx, constants.ERROR_LEVEL, attuneBookingResponse.Response.Message, nil, nil)
		attuneClient.Sentry.LogError(ctx, attuneBookingResponse.Response.Message, nil, map[string]interface{}{
			"code":     attuneBookingResponse.Response.Code,
			"message":  attuneBookingResponse.Response.Message,
			"visit_id": payload.OrderID,
		})
		return &structures.CommonError{
			Message:    attuneBookingResponse.Response.Message,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}

func (attuneClient *AttuneClient) SyncDataToAttuneAfterSync(ctx context.Context, labId uint,
	payload structures.AttuneSyncDataToLisAfterSyncRequest) *structures.CommonError {

	token := attuneClient.getToken(ctx)
	if token == "" {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_SENDING_ATTUNE_WEBHOOK, nil, nil)
		return &structures.CommonError{
			Message:    constants.ERROR_WHILE_SENDING_ATTUNE_WEBHOOK,
			StatusCode: http.StatusInternalServerError,
		}
	}

	orgCode := utils.GetAttuneOrgCodeByLabId(labId)
	if orgCode == "" {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_SENDING_ATTUNE_WEBHOOK, nil, nil)
		return &structures.CommonError{
			Message:    constants.ERROR_GETTING_ATTUNE_ORG_CODE,
			StatusCode: http.StatusInternalServerError,
		}
	}

	headers := attuneClient.Headers
	headers["Authorization"] = fmt.Sprintf("Bearer %s", token)
	headers["IsWrapReq"] = "Y"
	headers["OrgCode"] = orgCode
	headers["Content-Type"] = constants.ContentTypeJson

	utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), map[string]interface{}{
		"payload": payload,
	}, nil)

	var response interface{}
	err := attuneClient.ApiClient.Post(ctx, &response, ATTUNE_URLS[SYNC_DATA_TO_ATTUNE], nil, payload, headers,
		SyncDataToAttuneRetries, SyncDataToAttuneDelay)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_SENDING_ATTUNE_WEBHOOK, nil, err)
		return &structures.CommonError{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_SENDING_ATTUNE_WEBHOOK, nil, err)
		return &structures.CommonError{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	attuneBookingResponse := structures.AttuneBookingResponse{}
	err = json.Unmarshal([]byte(responseBytes), &attuneBookingResponse)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_SENDING_ATTUNE_WEBHOOK, nil, err)
		return &structures.CommonError{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	if strings.ToLower(attuneBookingResponse.Response.Code) != "success" {
		utils.AddLog(ctx, constants.ERROR_LEVEL, attuneBookingResponse.Response.Message, nil, nil)
		attuneClient.Sentry.LogError(ctx, attuneBookingResponse.Response.Message, nil, map[string]interface{}{
			"code":     attuneBookingResponse.Response.Code,
			"message":  attuneBookingResponse.Response.Message,
			"visit_id": payload.OrderID,
		})
		return &structures.CommonError{
			Message:    attuneBookingResponse.Response.Message,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}
