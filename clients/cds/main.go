package cdsClient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/clients"
	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/structures"
	"github.com/Orange-Health/citadel/common/utils"
)

type CdsClient struct {
	ApiClient *clients.ApiClient
	Headers   map[string]string
	Cache     cache.CacheLayer
}

func NewClient() *CdsClient {
	apiClient := clients.NewClient()
	apiClient.BaseURL = CdsBaseUrl
	return &CdsClient{
		ApiClient: apiClient,
		Headers: map[string]string{
			"service":     constants.CitadelServiceName,
			"api-key":     CdsApiKey,
			"retool-user": CdsUserEmail,
		},
		Cache: cache.InitializeCache(),
	}
}

func InitializeCdsClient() CdsClientInterface {
	return NewClient()
}

type CdsClientInterface interface {
	GetPanelDetails(ctx context.Context, testCodes []string, testIds []uint, cityCode string,
		labId uint, isTestCodes, isTestIds bool, patientDob, patientGender string) (
		structures.MasterInvestigationResponse, error)
	GetInvestigationDetails(ctx context.Context, testCodes []string, cityCode string, labId uint,
		patientDob, patientGender string) ([]structures.Investigation, error)
	GetDepartmentMapping(ctx context.Context) (map[string]string, error)
	GetDependentInvestigations(ctx context.Context, masterInvestigationId, labId uint,
		cityCode, dob, gender string) (structures.DependentInvestigationsApiResponse, error)
	GetMasterVialTypes(ctx context.Context) ([]structures.MasterVialType, error)
	GetLabMasters(ctx context.Context) ([]structures.Lab, error)
	GetCollectionSequences(ctx context.Context, collectionSequenceRequest structures.CollectionSequenceRequest) (
		structures.CollectionSequenceResponse, error)
	GetMasterTests(ctx context.Context, withTestInfo, withSources, withAliases, withVialTypeMaps,
		withTestCityMeta, withLabMeta bool, masterTestIds []uint, onlyActive bool) ([]structures.CdsTestMaster, error)
	GetDeduplicatedTestsAndPackages(ctx context.Context, testIds, packageIds []uint) (structures.TestDeduplicationResponse,
		error)
	GetNrlCpEnabledMasterTestIds(ctx context.Context) ([]uint, error)
}

func (cdsClient *CdsClient) GetPanelDetails(ctx context.Context, testCodes []string, testIds []uint, cityCode string,
	labId uint, isTestCodes, isTestIds bool, patientDob, patientGender string) (
	structures.MasterInvestigationResponse, error) {

	masterInvestigationResponse := structures.MasterInvestigationResponse{}
	var response interface{}

	headers := cdsClient.Headers
	headers["Content-Type"] = constants.ContentTypeJsonWithCharsetUtf8
	headers["city-code"] = cityCode

	payload := make(map[string]interface{})
	payload["lab_id"] = labId
	payload["dob"] = patientDob
	payload["gender"] = patientGender

	queryParams := make(map[string]interface{})
	queryParams["is_fetch_mrr"] = true

	if isTestCodes {
		payload["test_codes"] = testCodes
		queryParams["is_test_codes"] = true
	} else if isTestIds {
		payload["test_ids"] = testIds
		queryParams["is_test_ids"] = true
	}

	err := cdsClient.ApiClient.Post(ctx, &response, URL_MAP[PANEL_DETAILS], queryParams, payload, headers, 3,
		time.Millisecond*500)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_PANEL_DETAILS, nil, err)
		return masterInvestigationResponse, err
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_PANEL_DETAILS, nil, err)
		return masterInvestigationResponse, err
	}

	err = json.Unmarshal(responseBytes, &masterInvestigationResponse)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_PANEL_DETAILS, nil, err)
		return masterInvestigationResponse, err
	}

	return masterInvestigationResponse, nil
}

func (cdsClient *CdsClient) GetInvestigationDetails(ctx context.Context, testCodes []string, cityCode string, labId uint,
	patientDob, patientGender string) ([]structures.Investigation, error) {
	masterInvestigationResponse := structures.MasterInvestigationResponse{}
	var response interface{}
	var cacheResponse interface{}

	sort.Strings(testCodes)
	cacheKey := fmt.Sprintf(constants.CacheKeyInvestigationDetails, cityCode, utils.ConvertStringSliceToString(testCodes),
		labId, patientDob, patientGender)

	cacheErr := cdsClient.Cache.Get(ctx, cacheKey, &cacheResponse)
	if cacheErr == nil {
		responseBytes, cacheErr := json.Marshal(cacheResponse)
		if cacheErr == nil {
			investigations := []structures.Investigation{}
			cacheErr := json.Unmarshal(responseBytes, &investigations)
			if cacheErr == nil {
				return investigations, nil
			}
		}
	}

	headers := cdsClient.Headers
	headers["Content-Type"] = constants.ContentTypeJsonWithCharsetUtf8
	headers["city-code"] = cityCode

	payload := make(map[string]interface{})
	payload["lab_id"] = labId
	payload["dob"] = patientDob
	payload["gender"] = patientGender
	payload["test_codes"] = testCodes

	queryParams := make(map[string]interface{})
	queryParams["is_fetch_mrr"] = true

	err := cdsClient.ApiClient.Post(ctx, &response, URL_MAP[INVESTIGATION_DETAILS], queryParams, payload, headers, 3,
		time.Millisecond*500)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_INVESTIGATION_DETAILS, nil, err)
		return nil, err
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_INVESTIGATION_DETAILS, nil, err)
		return nil, err
	}

	err = json.Unmarshal(responseBytes, &masterInvestigationResponse)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_INVESTIGATION_DETAILS, nil, err)
		return nil, err
	}

	if len(masterInvestigationResponse.InvestigationDetails.Investigations) == 0 {
		return nil, errors.New(constants.ERROR_IN_GETTING_INVESTIGATION_DETAILS)
	}

	err = cdsClient.Cache.Set(ctx, cacheKey,
		masterInvestigationResponse.InvestigationDetails.Investigations[0], constants.CacheExpiry2MinutesInt)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_INVESTIGATION_DETAILS, nil, err)
	}

	return masterInvestigationResponse.InvestigationDetails.Investigations, nil
}

func (cdsClient *CdsClient) GetDepartmentMapping(ctx context.Context) (map[string]string, error) {

	departmentMappingMap, err := cdsClient.Cache.HGetAll(ctx, constants.CacheKeyDepartmentMapping)
	if err == nil && len(departmentMappingMap) > 0 {
		departmentMapping := make(map[string]string)
		for k, v := range departmentMappingMap {
			departmentMapping[k] = fmt.Sprint(v)
		}
		return departmentMapping, nil
	}

	var response interface{}
	headers := cdsClient.Headers
	headers["Content-Type"] = constants.ContentTypeJson
	err = cdsClient.ApiClient.Get(ctx, &response, URL_MAP[DEPARTMENT_MAPPING], nil, nil, headers, 3, time.Millisecond*500)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_DEPARTMENT_MAPPING, nil, err)
		return map[string]string{}, err
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_DEPARTMENT_MAPPING, nil, err)
		return map[string]string{}, err
	}

	var responseArray []map[string]interface{}
	err = json.Unmarshal(responseBytes, &responseArray)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_DEPARTMENT_MAPPING, nil, err)
		return map[string]string{}, err
	}

	departmentMapping := make(map[string]string)

	for _, mapping := range responseArray {
		department := mapping["department"].(string)
		masterTestId := mapping["master_test_id"]
		departmentMapping[fmt.Sprint(masterTestId)] = utils.ConvertStringToCamelCase(department)
	}

	err = cdsClient.Cache.HSetAll(ctx, constants.CacheKeyDepartmentMapping, departmentMapping, constants.CacheExpiry24Hours)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_DEPARTMENT_MAPPING, nil, err)
	}

	return departmentMapping, nil
}

func (cdsClient *CdsClient) GetDependentInvestigations(ctx context.Context, masterInvestigationId, labId uint,
	cityCode, dob, gender string) (structures.DependentInvestigationsApiResponse, error) {
	dependentInvestigations := structures.DependentInvestigationsApiResponse{}
	var response interface{}
	var cacheResponse interface{}
	cacheKey := fmt.Sprintf(constants.CacheKeyDependentInvestigations, cityCode, masterInvestigationId,
		labId, dob, gender)

	cacheErr := cdsClient.Cache.Get(ctx, cacheKey, &cacheResponse)
	if cacheErr == nil {
		responseBytes, cacheErr := json.Marshal(cacheResponse)
		if cacheErr == nil {
			cacheErr := json.Unmarshal(responseBytes, &dependentInvestigations)
			if cacheErr == nil {
				return dependentInvestigations, nil
			}
		}
	}

	headers := cdsClient.Headers
	headers["Content-Type"] = constants.ContentTypeJsonWithCharsetUtf8
	headers["city-code"] = cityCode

	queryParams := make(map[string]interface{})
	queryParams["investigation_id"] = masterInvestigationId
	queryParams["dob"] = dob
	queryParams["lab_id"] = labId
	queryParams["gender"] = gender

	err := cdsClient.ApiClient.Get(ctx, &response, URL_MAP[DEPENDENT_INVESTIGATIONS], queryParams, nil, headers, 3,
		time.Millisecond*500)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_DEPENDENT_INVESTIGATIONS, nil, err)
		return dependentInvestigations, err
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_DEPENDENT_INVESTIGATIONS, nil, err)
		return dependentInvestigations, err
	}

	err = json.Unmarshal(responseBytes, &dependentInvestigations)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_DEPENDENT_INVESTIGATIONS, nil, err)
		return dependentInvestigations, err
	}

	err = cdsClient.Cache.Set(ctx, cacheKey, dependentInvestigations, constants.CacheExpiry2MinutesInt)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_DEPENDENT_INVESTIGATIONS, nil, err)
	}

	return dependentInvestigations, nil
}
func (cdsClient *CdsClient) GetMasterVialTypes(ctx context.Context) ([]structures.MasterVialType, error) {
	var response interface{}

	headers := cdsClient.Headers
	headers["Content-Type"] = constants.ContentTypeJsonWithCharsetUtf8

	vialsResponse := VialsApiResponse{}
	err := cdsClient.ApiClient.Get(ctx, &response, URL_MAP[VIALS], nil, nil, headers, 3, time.Millisecond*500)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_VIALS, nil, err)
		return nil, err
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_VIALS, nil, err)
		return nil, err
	}

	err = json.Unmarshal(responseBytes, &vialsResponse)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_VIALS, nil, err)
		return nil, err
	}

	return vialsResponse.MasterVialTypes, nil
}

func (cdsClient *CdsClient) GetLabMasters(ctx context.Context) ([]structures.Lab, error) {
	var response interface{}

	headers := cdsClient.Headers
	headers["Content-Type"] = constants.ContentTypeJsonWithCharsetUtf8

	labsResponse := LabsApiResponse{}
	err := cdsClient.ApiClient.Get(ctx, &response, URL_MAP[LAB_MASTERS], nil, nil, headers, 3, time.Millisecond*500)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_LABS, nil, err)
		return nil, err
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_LABS, nil, err)
		return nil, err
	}

	err = json.Unmarshal(responseBytes, &labsResponse)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_LABS, nil, err)
		return nil, err
	}

	return labsResponse.Labs, nil
}

func (cdsClient *CdsClient) GetDeduplicatedTestsAndPackages(ctx context.Context, testIds, packageIds []uint) (
	structures.TestDeduplicationResponse, error) {
	var response interface{}
	if len(testIds) == 0 && len(packageIds) == 0 {
		return structures.TestDeduplicationResponse{}, errors.New(constants.ERROR_TESTS_AND_PACKAGES_CANNOT_BE_EMPTY)
	}

	headers := cdsClient.Headers
	headers["Content-Type"] = constants.ContentTypeJsonWithCharsetUtf8

	payload := make(map[string]interface{})
	payload["test_ids"] = testIds
	payload["package_ids"] = packageIds

	err := cdsClient.ApiClient.Post(ctx, &response, URL_MAP[TEST_DEDUPLICATION], nil, payload, headers, 3,
		time.Millisecond*500)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_DEDUPLICATION_RESPONSE, nil, err)
		return structures.TestDeduplicationResponse{}, err
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_DEDUPLICATION_RESPONSE, nil, err)
		return structures.TestDeduplicationResponse{}, err
	}

	testDeduplicationResponseData := structures.TestDeduplicationResponseData{}
	err = json.Unmarshal(responseBytes, &testDeduplicationResponseData)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_DEDUPLICATION_RESPONSE, nil, err)
		return structures.TestDeduplicationResponse{}, err
	}

	return testDeduplicationResponseData.Data, nil

}

func (cdsClient *CdsClient) GetMasterTests(ctx context.Context, withTestInfo, withSources, withAliases, withVialTypeMaps,
	withTestCityMeta, withLabMeta bool, masterTestIds []uint, onlyActive bool) ([]structures.CdsTestMaster, error) {
	var response interface{}
	headers := cdsClient.Headers
	headers["Content-Type"] = constants.ContentTypeJsonWithCharsetUtf8

	queryParams := map[string]interface{}{}
	if withTestInfo {
		queryParams["with_test_info"] = "true"
	} else {
		queryParams["with_test_info"] = "false"
	}
	if withSources {
		queryParams["with_sources"] = "true"
	}
	if withAliases {
		queryParams["with_aliases"] = "true"
	}
	if withVialTypeMaps {
		queryParams["with_vial_type_maps"] = "true"
	}
	if withTestCityMeta {
		queryParams["with_test_city_meta_for"] = strings.Join(constants.ActiveCityCodes, ",")
	}
	if withLabMeta {
		queryParams["with_lab_meta"] = "true"
	}
	if len(masterTestIds) > 0 {
		queryParams["q[ids]"] = utils.ConvertUintSliceToString(masterTestIds)
	}
	if onlyActive {
		queryParams["q[status]"] = "1"
	}
	queryParams["limit"] = "3000"

	masterTestsResponse := MasterTestsResponse{}
	err := cdsClient.ApiClient.Get(ctx, &response, URL_MAP[MASTER_TESTS], queryParams, nil, headers, 3, time.Millisecond*500)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_MASTER_TESTS, nil, err)
		return []structures.CdsTestMaster{}, err
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_MASTER_TESTS, nil, err)
		return nil, err
	}

	err = json.Unmarshal(responseBytes, &masterTestsResponse)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_MASTER_TESTS, nil, err)
		return nil, err
	}

	return masterTestsResponse.Tests, nil
}

func (cdsClient *CdsClient) GetCollectionSequences(ctx context.Context,
	collectionSequenceRequest structures.CollectionSequenceRequest) (structures.CollectionSequenceResponse, error) {
	var response interface{}
	collectionSequenceResponse := structures.CollectionSequenceResponse{}

	headers := cdsClient.Headers
	headers["Content-Type"] = constants.ContentTypeJsonWithCharsetUtf8

	err := cdsClient.ApiClient.Post(ctx, &response, URL_MAP[COLLECTION_SEQUENCE], nil, collectionSequenceRequest, headers, 3,
		time.Millisecond*500)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_COLLECTION_SEQUENCE, nil, err)
		return collectionSequenceResponse, err
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_COLLECTION_SEQUENCE, nil, err)
		return collectionSequenceResponse, err
	}

	err = json.Unmarshal(responseBytes, &collectionSequenceResponse)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_COLLECTION_SEQUENCE, nil, err)
		return collectionSequenceResponse, err
	}

	return collectionSequenceResponse, nil
}

func (cdsClient *CdsClient) GetNrlCpEnabledMasterTestIds(ctx context.Context) ([]uint, error) {
	var response interface{}
	nrlCpEnabledMasterTestIds := []uint{}

	headers := cdsClient.Headers
	headers["Content-Type"] = constants.ContentTypeJsonWithCharsetUtf8

	err := cdsClient.ApiClient.Get(ctx, &response, URL_MAP[NRL_ENABLED_TESTS], nil, nil, headers, 3, time.Millisecond*500)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_NRL_ENABLED_TESTS, nil, err)
		return nrlCpEnabledMasterTestIds, err
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_NRL_ENABLED_TESTS, nil, err)
		return nrlCpEnabledMasterTestIds, err
	}

	nrlEnabledMasterTestIdResponse := NrlEnabledMasterTestIdsResponse{}
	err = json.Unmarshal(responseBytes, &nrlEnabledMasterTestIdResponse)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_NRL_ENABLED_TESTS, nil, err)
		return nrlCpEnabledMasterTestIds, err
	}

	for _, test := range nrlEnabledMasterTestIdResponse.Data {
		nrlCpEnabledMasterTestIds = append(nrlCpEnabledMasterTestIds, test.Id)
	}

	return nrlCpEnabledMasterTestIds, nil
}
