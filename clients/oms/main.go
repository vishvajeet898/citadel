package omsClient

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/Orange-Health/citadel/clients"
	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/structures"
	"github.com/Orange-Health/citadel/common/utils"
)

type OmsClient struct {
	ApiClient *clients.ApiClient
	Headers   map[string]string
}

func NewClient() *OmsClient {
	apiClient := clients.NewClient()
	apiClient.BaseURL = OmsBaseUrl
	return &OmsClient{
		ApiClient: apiClient,
		Headers: map[string]string{
			"service": constants.CitadelServiceName,
			"api-key": OmsApiKey,
		},
	}
}

func InitializeOmsClient() OmsClientInterface {
	return NewClient()
}

type OmsClientInterface interface {
	GetDeltaValuesFromPatientIds(ctx context.Context, patientIds []string, masterInvestigationIds []uint) (
		[]structures.DeltaValuesStruct, error)
	GetPatientPastRecordsFromPatientIds(ctx context.Context, patientIds []string) (
		[]structures.PatientPastRecords, error)
	GetOrderById(ctx context.Context, omsOrderId, cityCode string) (structures.OmsOrderDetailsResponse, error)

	UpdatePatientDetails(ctx context.Context, orderId uint, cityCode string, patientDetails interface{}) (interface{}, error)
}

func (omsClient *OmsClient) UpdatePatientDetails(ctx context.Context, orderId uint, cityCode string,
	patientDetails interface{}) (interface{}, error) {
	var response interface{}

	headers := omsClient.Headers
	headers["Content-Type"] = constants.ContentTypeJsonWithCharsetUtf8

	url := URL_MAP[UPDATE_PATIENT_DETAILS]
	url = fmt.Sprintf(url, orderId)

	queryParams := map[string]interface{}{
		"city_code": cityCode,
	}

	err := omsClient.ApiClient.Put(ctx, &response, url, queryParams, patientDetails, headers, 3, time.Millisecond*500)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_UPDATING_PATIENT_DETAILS, nil, err)
		return nil, err
	}

	return response, nil
}

func (omsClient *OmsClient) GetDeltaValuesFromPatientIds(ctx context.Context, patientIds []string,
	masterInvestigationIds []uint) ([]structures.DeltaValuesStruct, error) {
	deltaValuesStructs := []structures.DeltaValuesStruct{}
	cityCodes := constants.ActiveCityCodes
	deltaValueResponses := make([]structures.DeltaValuesResponse, len(cityCodes))
	var response interface{}
	headers := omsClient.Headers
	headers["Content-Type"] = constants.ContentTypeJsonWithCharsetUtf8

	url := URL_MAP[DELTA_VALUES_FROM_PATIENT_IDS]
	patientIdsString := utils.ConvertStringSliceToString(patientIds)
	masterInvestigationIdsString := utils.ConvertUintSliceToString(masterInvestigationIds)

	wg := sync.WaitGroup{}

	for index, cityCode := range cityCodes {
		wg.Add(1)

		go func(index int, cityCode, patientIdsString, masterInvestigationIdsString string, headers map[string]string) {
			defer wg.Done()

			queryParams := map[string]interface{}{
				"patient_ids":              patientIdsString,
				"master_investigation_ids": masterInvestigationIdsString,
				"city_code":                cityCode,
			}
			err := omsClient.ApiClient.Get(ctx, &response, url, queryParams, nil, headers, 1, 0)
			if err != nil {
				utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_FETCHING_DELTA_VALUES, nil, err)
				return
			}

			deltaValueResponse := structures.DeltaValuesResponse{}
			responseBytes, err := json.Marshal(response)
			if err != nil {
				utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_FETCHING_DELTA_VALUES, nil, err)
				return
			}

			err = json.Unmarshal(responseBytes, &deltaValueResponse)
			if err != nil {
				utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_FETCHING_DELTA_VALUES, nil, err)
				return
			}

			for i := range deltaValueResponse.Data {
				deltaValueResponse.Data[i].CityCode = cityCode
			}

			deltaValueResponses[index] = deltaValueResponse
		}(index, cityCode, patientIdsString, masterInvestigationIdsString, headers)
	}

	wg.Wait()

	for _, deltaValueResponse := range deltaValueResponses {
		if len(deltaValueResponse.Data) != 0 {
			deltaValuesStructs = append(deltaValuesStructs, deltaValueResponse.Data...)
		}
	}

	return deltaValuesStructs, nil
}

func (omsClient *OmsClient) GetPatientPastRecordsFromPatientIds(ctx context.Context, patientIds []string) (
	[]structures.PatientPastRecords, error) {
	patientPastRecords := []structures.PatientPastRecords{}
	cityCodes := constants.ActiveCityCodes
	patientPastRecordsResponses := make([]structures.PatientPastRecordsResponse, len(cityCodes))
	var response interface{}
	headers := omsClient.Headers
	headers["Content-Type"] = constants.ContentTypeJsonWithCharsetUtf8

	url := URL_MAP[PATIENT_PAST_RECORDS_FROM_PATIENT_IDS]
	patientIdsString := utils.ConvertStringSliceToString(patientIds)

	wg := sync.WaitGroup{}

	for index, cityCode := range cityCodes {
		wg.Add(1)

		go func(index int, cityCode, patientIdsString string, headers map[string]string) {
			defer wg.Done()

			queryParams := map[string]interface{}{
				"patient_ids": patientIdsString,
				"city_code":   cityCode,
			}
			err := omsClient.ApiClient.Get(ctx, &response, url, queryParams, nil, headers, 1, 0)
			if err != nil {
				utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_FETCHING_PATIENT_PAST_RECORDS,
					nil, err)
				return
			}

			responseBytes, err := json.Marshal(response)
			if err != nil {
				utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_FETCHING_PATIENT_PAST_RECORDS,
					nil, err)
				return
			}

			patientPastRecordResponse := structures.PatientPastRecordsResponse{}
			err = json.Unmarshal(responseBytes, &patientPastRecordResponse)
			if err != nil {
				utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_FETCHING_PATIENT_PAST_RECORDS,
					nil, err)
				return
			}

			for i := range patientPastRecordResponse.Data {
				patientPastRecordResponse.Data[i].CityCode = cityCode
			}

			patientPastRecordsResponses[index] = patientPastRecordResponse
		}(index, cityCode, patientIdsString, headers)
	}

	wg.Wait()

	for _, patientPastRecordsResponse := range patientPastRecordsResponses {
		if len(patientPastRecordsResponse.Data) != 0 {
			patientPastRecords = append(patientPastRecords, patientPastRecordsResponse.Data...)
		}
	}

	return patientPastRecords, nil
}

func (omsClient *OmsClient) GetOrderById(ctx context.Context, omsOrderId, cityCode string) (structures.OmsOrderDetailsResponse, error) {
	var response interface{}
	headers := omsClient.Headers
	headers["Content-Type"] = constants.ContentTypeJsonWithCharsetUtf8

	url := fmt.Sprintf(URL_MAP[GET_ORDER_BY_ID], omsOrderId)

	queryParams := map[string]interface{}{
		"city_code": cityCode,
	}

	err := omsClient.ApiClient.Get(ctx, &response, url, queryParams, nil, headers, 1, 0)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_FETCHING_ORDER_DETAILS, nil, err)
		return structures.OmsOrderDetailsResponse{}, err
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_FETCHING_ORDER_DETAILS, nil, err)
		return structures.OmsOrderDetailsResponse{}, err
	}

	orderDetailsResponse := structures.OmsOrderDetailsResponse{}
	err = json.Unmarshal(responseBytes, &orderDetailsResponse)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_FETCHING_ORDER_DETAILS, nil, err)
		return structures.OmsOrderDetailsResponse{}, err
	}

	return orderDetailsResponse, nil
}
