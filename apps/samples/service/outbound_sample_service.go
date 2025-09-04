package service

import (
	"context"
	"net/http"

	"gorm.io/gorm"

	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

func (sampleService *SampleService) CreateInterlabSamplesWithTx(ctx context.Context, tx *gorm.DB,
	samples []commonModels.Sample, samplesMetadata []commonModels.SampleMetadata,
	sampleNumberToLabIdMap map[uint]uint) *commonStructures.CommonError {
	if len(samples) == 0 {
		return &commonStructures.CommonError{
			StatusCode: http.StatusBadRequest,
			Message:    commonConstants.ERROR_NO_INTERLAB_SAMPLES,
		}
	}

	sampleIdToSampleMetadataMap := map[uint]commonModels.SampleMetadata{}
	newSamples, newSamplesMetadata := []commonModels.Sample{}, []commonModels.SampleMetadata{}

	for index := range samplesMetadata {
		sampleIdToSampleMetadataMap[samplesMetadata[index].SampleId] = samplesMetadata[index]
	}

	for _, sample := range samples {
		newSample := commonModels.Sample{
			OmsOrderId:       sample.OmsOrderId,
			OmsRequestId:     sample.OmsRequestId,
			OmsCityCode:      sample.OmsCityCode,
			VialTypeId:       sample.VialTypeId,
			Barcode:          sample.Barcode,
			Status:           commonConstants.SampleInTransfer,
			LabId:            0,
			DestinationLabId: sampleNumberToLabIdMap[sample.SampleNumber],
			ParentSampleId:   sample.Id,
			SampleNumber:     sample.SampleNumber,
		}
		newSample.CreatedBy = commonConstants.CitadelSystemId
		newSample.UpdatedBy = commonConstants.CitadelSystemId
		currentTime := commonUtils.GetCurrentTime()
		newSampleMetadata := commonModels.SampleMetadata{
			OmsCityCode:              sampleIdToSampleMetadataMap[sample.Id].OmsCityCode,
			OmsOrderId:               sampleIdToSampleMetadataMap[sample.Id].OmsOrderId,
			CollectionSequenceNumber: sampleIdToSampleMetadataMap[sample.Id].CollectionSequenceNumber,
			LastUpdatedAt:            currentTime,
			TransferredAt:            currentTime,
			CollectedAt:              sampleIdToSampleMetadataMap[sample.Id].CollectedAt,
			ReceivedAt:               sampleIdToSampleMetadataMap[sample.Id].ReceivedAt,
			RejectedAt:               sampleIdToSampleMetadataMap[sample.Id].RejectedAt,
			NotReceivedAt:            sampleIdToSampleMetadataMap[sample.Id].NotReceivedAt,
			LisSyncAt:                sampleIdToSampleMetadataMap[sample.Id].LisSyncAt,
			BarcodeScannedAt:         sampleIdToSampleMetadataMap[sample.Id].BarcodeScannedAt,
			NotReceivedReason:        sampleIdToSampleMetadataMap[sample.Id].NotReceivedReason,
			BarcodeImageUrl:          sampleIdToSampleMetadataMap[sample.Id].BarcodeImageUrl,
			TaskSequence:             sampleIdToSampleMetadataMap[sample.Id].TaskSequence,
			CollectLaterReason:       sampleIdToSampleMetadataMap[sample.Id].CollectLaterReason,
			RejectingLab:             sampleIdToSampleMetadataMap[sample.Id].RejectingLab,
			CollectedVolume:          sampleIdToSampleMetadataMap[sample.Id].CollectedVolume,
		}
		newSampleMetadata.CreatedBy = commonConstants.CitadelSystemId
		newSampleMetadata.UpdatedBy = commonConstants.CitadelSystemId
		newSamples = append(newSamples, newSample)
		newSamplesMetadata = append(newSamplesMetadata, newSampleMetadata)
	}

	if cErr := sampleService.SampleDao.CreateInterlabSamplesWithTx(ctx, tx, newSamples, newSamplesMetadata); cErr != nil {
		return cErr
	}

	return nil
}

func (sampleService *SampleService) GetInterlabSamplesByParentSampleId(ctx context.Context, parentSampleId uint) ([]commonStructures.SampleInfo, *commonStructures.CommonError) {
	samples, cErr := sampleService.SampleDao.GetInterlabSamplesByParentSampleId(ctx, parentSampleId)
	if cErr != nil {
		return nil, cErr
	}

	return samples, nil
}
