package mapper

import (
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

func MapSampleSampleMetaToSampleInfo(sample commonModels.Sample,
	sampleMeta commonModels.SampleMetadata) commonStructures.SampleInfo {
	return commonStructures.SampleInfo{
		Id:                       sample.Id,
		OmsCityCode:              sample.OmsCityCode,
		OmsOrderId:               sample.OmsOrderId,
		OmsRequestId:             sample.OmsRequestId,
		VisitId:                  sample.VisitId,
		VialTypeId:               sample.VialTypeId,
		Barcode:                  sample.Barcode,
		Status:                   sample.Status,
		LabId:                    sample.LabId,
		DestinationLabId:         sample.DestinationLabId,
		RejectionReason:          sample.RejectionReason,
		ParentSampleId:           sample.ParentSampleId,
		SampleNumber:             sample.SampleNumber,
		SampleId:                 sampleMeta.SampleId,
		CollectionSequenceNumber: sampleMeta.CollectionSequenceNumber,
		LastUpdatedAt:            sampleMeta.LastUpdatedAt,
		TransferredAt:            sampleMeta.TransferredAt,
		OutsourcedAt:             sampleMeta.OutsourcedAt,
		CollectedAt:              sampleMeta.CollectedAt,
		ReceivedAt:               sampleMeta.ReceivedAt,
		RejectedAt:               sampleMeta.RejectedAt,
		NotReceivedAt:            sampleMeta.NotReceivedAt,
		LisSyncAt:                sampleMeta.LisSyncAt,
		BarcodeScannedAt:         sampleMeta.BarcodeScannedAt,
		NotReceivedReason:        sampleMeta.NotReceivedReason,
		BarcodeImageUrl:          sampleMeta.BarcodeImageUrl,
		TaskSequence:             sampleMeta.TaskSequence,
		CollectLaterReason:       sampleMeta.CollectLaterReason,
		RejectingLab:             sampleMeta.RejectingLab,
		CollectedVolume:          sampleMeta.CollectedVolume,
		CreatedAt:                sample.CreatedAt,
		CreatedBy:                sample.CreatedBy,
		UpdatedAt:                sample.UpdatedAt,
		UpdatedBy:                sample.UpdatedBy,
		DeletedBy:                sample.DeletedBy,
		DeletedAt:                commonUtils.GetGoLangTimeFromGormDeletedAt(sample.DeletedAt),
	}
}

func MapBulkSampleSampleMetaToSampleInfo(sample []commonModels.Sample,
	sampleMeta []commonModels.SampleMetadata) []commonStructures.SampleInfo {
	sampleMetaMap := make(map[uint]commonModels.SampleMetadata)
	for _, s := range sampleMeta {
		sampleMetaMap[s.SampleId] = s
	}

	var sampleInfos []commonStructures.SampleInfo
	for _, s := range sample {
		sampleInfos = append(sampleInfos, MapSampleSampleMetaToSampleInfo(s, sampleMetaMap[s.Id]))
	}
	return sampleInfos
}

func MapSampleInfoToSample(sampleInfo commonStructures.SampleInfo) commonModels.Sample {

	baseModel := commonModels.BaseModel{
		Id:        sampleInfo.Id,
		CreatedAt: sampleInfo.CreatedAt,
		CreatedBy: sampleInfo.CreatedBy,
		UpdatedAt: sampleInfo.UpdatedAt,
		UpdatedBy: sampleInfo.UpdatedBy,
		DeletedBy: sampleInfo.DeletedBy,
		DeletedAt: commonUtils.GetGormDeletedAtFromGoLangTime(sampleInfo.DeletedAt),
	}

	return commonModels.Sample{
		OmsCityCode:      sampleInfo.OmsCityCode,
		OmsOrderId:       sampleInfo.OmsOrderId,
		OmsRequestId:     sampleInfo.OmsRequestId,
		VisitId:          sampleInfo.VisitId,
		VialTypeId:       sampleInfo.VialTypeId,
		Barcode:          sampleInfo.Barcode,
		Status:           sampleInfo.Status,
		LabId:            sampleInfo.LabId,
		DestinationLabId: sampleInfo.DestinationLabId,
		RejectionReason:  sampleInfo.RejectionReason,
		ParentSampleId:   sampleInfo.ParentSampleId,
		SampleNumber:     sampleInfo.SampleNumber,
		BaseModel:        baseModel,
	}
}

func MapBulkSampleInfoToSample(sampleInfo []commonStructures.SampleInfo) []commonModels.Sample {
	var samples []commonModels.Sample
	for _, s := range sampleInfo {
		samples = append(samples, MapSampleInfoToSample(s))
	}
	return samples
}

func MapSampleInfoToSampleMetadata(sampleInfo commonStructures.SampleInfo) commonModels.SampleMetadata {
	baseModel := commonModels.BaseModel{
		CreatedAt: sampleInfo.CreatedAt,
		CreatedBy: sampleInfo.CreatedBy,
		UpdatedAt: sampleInfo.UpdatedAt,
		UpdatedBy: sampleInfo.UpdatedBy,
		DeletedBy: sampleInfo.DeletedBy,
		DeletedAt: commonUtils.GetGormDeletedAtFromGoLangTime(sampleInfo.DeletedAt),
	}
	return commonModels.SampleMetadata{
		OmsCityCode:              sampleInfo.OmsCityCode,
		OmsOrderId:               sampleInfo.OmsOrderId,
		SampleId:                 sampleInfo.SampleId,
		CollectionSequenceNumber: sampleInfo.CollectionSequenceNumber,
		BarcodeImageUrl:          sampleInfo.BarcodeImageUrl,
		LastUpdatedAt:            sampleInfo.LastUpdatedAt,
		LisSyncAt:                sampleInfo.LisSyncAt,
		BarcodeScannedAt:         sampleInfo.BarcodeScannedAt,
		CollectedAt:              sampleInfo.CollectedAt,
		ReceivedAt:               sampleInfo.ReceivedAt,
		RejectedAt:               sampleInfo.RejectedAt,
		NotReceivedAt:            sampleInfo.NotReceivedAt,
		TransferredAt:            sampleInfo.TransferredAt,
		OutsourcedAt:             sampleInfo.OutsourcedAt,
		RejectingLab:             sampleInfo.RejectingLab,
		CollectLaterReason:       sampleInfo.CollectLaterReason,
		TaskSequence:             sampleInfo.TaskSequence,
		CollectedVolume:          sampleInfo.CollectedVolume,
		BaseModel:                baseModel,
	}
}

func MapBulkSampleInfoToSampleMetadata(sampleInfo []commonStructures.SampleInfo) []commonModels.SampleMetadata {
	var sampleMetas []commonModels.SampleMetadata
	for _, s := range sampleInfo {
		sampleMetas = append(sampleMetas, MapSampleInfoToSampleMetadata(s))
	}
	return sampleMetas
}
