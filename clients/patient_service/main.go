package patientServiceClient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/clients"
	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/structures"
	"github.com/Orange-Health/citadel/common/utils"
)

type PatientServiceClient struct {
	ApiClient *clients.ApiClient
	Headers   map[string]string
	Cache     cache.CacheLayer
}

func NewClient() *PatientServiceClient {
	apiClient := clients.NewClient()
	apiClient.BaseURL = PatientServiceBaseUrl
	return &PatientServiceClient{
		ApiClient: apiClient,
		Headers: map[string]string{
			"service": constants.CitadelServiceName,
			"api-key": PatientServiceApiKey,
		},
		Cache: cache.InitializeCache(),
	}
}

func InitializePatientServiceClient() PatientServiceClientInterface {
	return NewClient()
}

type PatientServiceClientInterface interface {
	GetPatientDetailsByPatientId(ctx context.Context, patientID string) (structures.PatientDetailsResponse, error)
	GetSimilarPatientDetails(ctx context.Context, patientID string) ([]structures.PatientDetailsResponse, error)
}

func (patientServiceClient *PatientServiceClient) GetPatientDetailsByPatientId(ctx context.Context, patientID string,
) (structures.PatientDetailsResponse, error) {
	var response interface{}
	patientDetailsResponse := []structures.PatientDetailsResponse{}

	headers := patientServiceClient.Headers
	headers["Content-Type"] = constants.ContentTypeJson

	queryParams := map[string]interface{}{
		"patientID": patientID,
	}

	err := patientServiceClient.ApiClient.Get(ctx, &response, URL_MAP[PATIENT_DETAILS_API], queryParams, nil, headers, 3, time.Millisecond*500)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_GETTING_PATIENT_DETAILS, nil, err)
		return structures.PatientDetailsResponse{}, err
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_GETTING_PATIENT_DETAILS, nil, err)
		return structures.PatientDetailsResponse{}, err
	}

	err = json.Unmarshal(responseBytes, &patientDetailsResponse)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_GETTING_PATIENT_DETAILS, nil, err)
		return structures.PatientDetailsResponse{}, err
	}

	if len(patientDetailsResponse) == 0 {
		return structures.PatientDetailsResponse{}, errors.New(constants.ERROR_WHILE_GETTING_PATIENT_DETAILS)
	}

	return patientDetailsResponse[0], nil
}

func (patientServiceClient *PatientServiceClient) GetSimilarPatientDetails(ctx context.Context, patientID string) (
	[]structures.PatientDetailsResponse, error,
) {
	var response interface{}
	similarPatientResponse := []structures.PatientDetailsResponse{}

	urlPath := fmt.Sprintf(URL_MAP[SIMILAR_PATIENTS_API], patientID)
	headers := patientServiceClient.Headers
	headers["Content-Type"] = constants.ContentTypeJson

	err := patientServiceClient.ApiClient.Get(ctx, &response, urlPath, nil, nil, headers, 1, 0)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_GETTING_SIMILAR_PATIENT_DETAILS, nil, err)
		return similarPatientResponse, err
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_GETTING_SIMILAR_PATIENT_DETAILS, nil, err)
		return similarPatientResponse, err
	}

	err = json.Unmarshal(responseBytes, &similarPatientResponse)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_GETTING_SIMILAR_PATIENT_DETAILS, nil, err)
		return similarPatientResponse, err
	}

	return similarPatientResponse, nil
}
