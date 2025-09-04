package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Orange-Health/citadel/apps/samples/structures"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
)

func (s *Sample) GetSampleDetailsByOrderId(c *gin.Context) {
	orderId := c.Param("orderId")

	if orderId == "" {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.ERROR_INVALID_ORDER_ID)
		return
	}

	sampleDetails, err := s.SampleService.GetSampleDetailsByOmsOrderIdAndLabId(orderId)
	if err != nil {
		c.JSON(http.StatusOK, commonStructures.CommonAPIResponse{
			Data: sampleDetails,
		})
		return
	}

	c.JSON(http.StatusOK, commonStructures.CommonAPIResponse{
		Data: sampleDetails,
	})
}

func (s *Sample) GetSampleDetailsByRequestId(c *gin.Context) {
	requestId := c.Param("requestId")

	if requestId == "" {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.ERROR_INVALID_REQUEST_ID)
		return
	}

	sampleDetails, cErr := s.SampleService.GetSampleDetailsByOmsRequestId(requestId)
	if cErr != nil {
		c.JSON(http.StatusOK, commonStructures.CommonAPIResponse{
			Data: sampleDetails,
		})
		return
	}
	c.JSON(http.StatusOK, commonStructures.CommonAPIResponse{
		Data: sampleDetails,
	})
}

func (s *Sample) GetSampleDetailsByOmsOrderIds(c *gin.Context) {
	orderIds := c.Query("order_ids")
	if orderIds == "" {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.ERROR_INVALID_ORDER_ID)
		return
	}
	orderIdsList := commonUtils.ConvertStringToStringSlice(orderIds)

	sampleDetails, err := s.SampleService.GetSampleDetailsByOmsOrderIds(orderIdsList)
	if err != nil {
		c.JSON(http.StatusOK, commonStructures.CommonAPIResponse{
			Data: sampleDetails,
		})
		return
	}

	c.JSON(http.StatusOK, commonStructures.CommonAPIResponse{
		Data: sampleDetails,
	})
}

func (s *Sample) GetLisSyncDataByVisitId(c *gin.Context) {
	visitId := c.Param("visitId")
	if visitId == "" {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.ERROR_VISIT_ID_NOT_FOUND)
		return
	}

	lisSyncDetails, err := s.SampleService.GetLisSyncDataByVisitId(c.Request.Context(), visitId)
	if err != nil {
		commonUtils.HandleErrorResponse(c, err.StatusCode, err.Message)
		return
	}

	c.JSON(http.StatusOK, commonStructures.CommonAPIResponse{
		Data: lisSyncDetails,
	})
}

func (s *Sample) UpdateSampleCollected(c *gin.Context) {
	userId, cErr := commonUtils.GetUserIdFromContext(c)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	sampleCollectionRequest := commonStructures.SampleCollectedRequest{}
	if err := c.ShouldBindJSON(&sampleCollectionRequest); err != nil {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	sampleCollectionRequest.UserId = userId

	err := s.SampleService.UpdateSampleCollected(sampleCollectionRequest)
	if err != nil {
		commonUtils.HandleErrorResponse(c, http.StatusInternalServerError, err.Message)
		return
	}
	c.JSON(http.StatusOK, commonStructures.CommonAPIResponse{
		Message: commonConstants.DONE_RESPONSE,
	})
}

func (s *Sample) AddBarcodeDetails(c *gin.Context) {
	barcodeDetails := structures.AddBarcodesRequest{}
	if err := c.ShouldBindJSON(&barcodeDetails); err != nil {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	testIdStatusMap, cErr := s.SampleService.AddBarcodeDetails(barcodeDetails)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}
	c.JSON(http.StatusOK, commonStructures.CommonAPIResponse{
		Data: testIdStatusMap,
	})
}

func (s *Sample) AddBarcodeDetailsForOrangers(c *gin.Context) {
	accessionBody := structures.UpdateAccessionBody{}
	if err := c.ShouldBindJSON(&accessionBody); err != nil {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	testIdStatusMap, cErr := s.SampleService.AddBarcodeDetailsForOrangers(accessionBody)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}
	c.JSON(http.StatusOK, commonStructures.CommonAPIResponse{
		Data: testIdStatusMap,
	})
}

func (s *Sample) GetSampleDetailsForScheduler(c *gin.Context) {
	sampleDetailsRequest := structures.SampleDetailsRequest{}
	if err := c.ShouldBindJSON(&sampleDetailsRequest); err != nil {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	sampleDetails, cErr := s.SampleService.GetSampleDetailsForScheduler(sampleDetailsRequest)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, http.StatusInternalServerError, cErr.Message)
		return
	}

	c.JSON(http.StatusOK, commonStructures.CommonAPIResponse{
		Data: sampleDetails,
	})
}

func (s *Sample) RejectSampleByBarcode(c *gin.Context) {
	barcode := c.Params.ByName("barcode")

	requestBody := structures.RejectSampleRequest{}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, commonStructures.CommonAPIResponse{Error: err.Error()})
		return
	}

	userId, cErr := commonUtils.GetUserIdFromContext(c)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}
	labId, _ := commonUtils.GetCurrentLabIdFromContext(c)
	requestBody.UserId = userId
	requestBody.LabId = labId

	omsOrderId, omsTestIds, cErr := s.SampleService.RejectSampleByBarcode(c.Request.Context(), barcode, requestBody)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	_ = s.SampleWorkerService.SampleRejectionTask(c.Request.Context(), omsOrderId, omsTestIds)

	c.JSON(http.StatusOK, commonStructures.CommonAPIResponse{Message: commonConstants.DONE_RESPONSE})
}

func (s *Sample) RejectSamplePartially(c *gin.Context) {
	requestBody := structures.RejectSamplePartiallyRequest{}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, commonStructures.CommonAPIResponse{Error: err.Error()})
		return
	}

	userId, cErr := commonUtils.GetUserIdFromContext(c)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}
	labId, _ := commonUtils.GetCurrentLabIdFromContext(c)
	requestBody.UserId = userId
	requestBody.LabId = labId

	omsOrderId, omsTestIds, cErr := s.SampleService.RejectSamplePartiallyBySampleNumberAndTestId(c.Request.Context(), requestBody)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	_ = s.SampleWorkerService.SampleRejectionTask(c.Request.Context(), omsOrderId, omsTestIds)

	c.JSON(http.StatusOK, commonStructures.CommonAPIResponse{Message: "Sample partially rejected"})
}

func (s *Sample) UpdateSampleDetailsPostTaskCompletion(c *gin.Context) {
	requestBody := commonStructures.UpdateSampleDetailsPostTaskCompletionRequest{}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := s.SampleService.UpdateSampleDetailsPostTaskCompletion(c.Request.Context(), requestBody)
	if err != nil {
		commonUtils.HandleErrorResponse(c, http.StatusInternalServerError, err.Message)
		return
	}

	c.JSON(http.StatusOK, commonStructures.CommonAPIResponse{
		Message: commonConstants.DONE_RESPONSE,
	})
}

func (s *Sample) UpdateTaskSequenceForSample(c *gin.Context) {
	requestBody := commonStructures.UpdateTaskSequenceRequest{}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	cErr := s.SampleService.UpdateTaskSequenceForSample(requestBody)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	c.JSON(http.StatusOK, commonStructures.CommonAPIResponse{
		Message: commonConstants.DONE_RESPONSE,
	})
}

func (s *Sample) RemapSamplesToNewTaskSequence(c *gin.Context) {
	requestBody := structures.RemapSamplesRequest{}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	cErr := s.SampleService.RemapSamplesToNewTaskSequence(requestBody)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	c.JSON(http.StatusOK, commonStructures.CommonAPIResponse{
		Message: commonConstants.DONE_RESPONSE,
	})
}

func (s *Sample) UpdateSampleDetailsForReschedule(c *gin.Context) {
	requestBody := structures.UpdateSampleDetailsForRescheduleRequest{}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	cErr := s.SampleService.UpdateSampleDetailsForReschedule(requestBody)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	c.JSON(http.StatusOK, commonStructures.CommonAPIResponse{
		Message: commonConstants.DONE_RESPONSE,
	})
}

func (s *Sample) GetOmsTestDetailsByVisitId(c *gin.Context) {
	visitId := c.Param("visitId")
	if visitId == "" {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.ERROR_VISIT_ID_NOT_FOUND)
		return
	}

	omsTestDetails, cErr := s.SampleService.GetOmsTestDetailsByVisitId(visitId)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
	}

	c.JSON(http.StatusOK, commonStructures.CommonAPIResponse{
		Data: omsTestDetails,
	})
}

func (s *Sample) CollectionPortalMarkAccessionAsAccessioned(c *gin.Context) {
	requestBody := structures.MarkAccessionAsAccessionedRequest{}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	cErr := s.SampleService.CollectionPortalMarkAccessionAsAccessioned(requestBody)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	c.JSON(http.StatusOK, commonStructures.CommonAPIResponse{
		Message: commonConstants.DONE_RESPONSE,
	})
}

func (s *Sample) AddCollectedVolumeToSample(c *gin.Context) {
	requestBody := structures.AddCollectedVolumneRequest{}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response, cErr := s.SampleService.AddCollectedVolumeToSample(requestBody)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	c.JSON(http.StatusOK, commonStructures.CommonAPIResponse{
		Data: response,
	})
}

func (s *Sample) ForcefullyMarkSampleAsCollected(c *gin.Context) {
	requestBody := structures.ForcefullyMarkSampleAsCollectedRequest{}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userId, cErr := commonUtils.GetUserIdFromContext(c)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}
	requestBody.UserId = userId

	cErr = s.SampleService.ForcefullyMarkSampleAsCollected(requestBody)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	c.JSON(http.StatusOK, commonStructures.CommonAPIResponse{
		Message: commonConstants.DONE_RESPONSE,
	})
}

func (s *Sample) GetSamplesForDelayedReverseLogisticsDashboard(c *gin.Context) {
	cityCode := c.Query("city_code")
	samples := s.SampleService.GetSamplesForDelayedReverseLogisticsDashboard(c.Request.Context(), cityCode)
	c.JSON(http.StatusOK, commonStructures.CommonAPIResponse{
		Data: samples,
	})
}

func (s *Sample) GetSrfOrderIds(c *gin.Context) {
	cityCode := c.Query("city_code")
	orderIds := s.SampleService.GetSrfOrderIds(c.Request.Context(), cityCode)
	c.JSON(http.StatusOK, commonStructures.CommonAPIResponse{
		Data: orderIds,
	})
}
