package dao

import (
	"context"

	"gorm.io/gorm"

	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

func (sampleDao *SampleDao) CreateInterlabSamplesWithTx(ctx context.Context, tx *gorm.DB, samples []commonModels.Sample,
	samplesMetadata []commonModels.SampleMetadata) *commonStructures.CommonError {

	for index := range samples {
		sample := &samples[index]
		sampleMetadata := &samplesMetadata[index]
		if err := tx.Create(&sample).Error; err != nil {
			return commonUtils.HandleORMError(err)
		}

		sampleMetadata.SampleId = sample.Id
		if err := tx.Create(&sampleMetadata).Error; err != nil {
			return commonUtils.HandleORMError(err)
		}
	}

	return nil
}

func (sampleDao *SampleDao) GetInterlabSamplesByParentSampleId(ctx context.Context, parentSampleId uint) (
	[]commonStructures.SampleInfo, *commonStructures.CommonError) {
	samples := []commonStructures.SampleInfo{}
	if err := sampleDao.Db.Table("samples").
		Select("sample_metadata.*, samples.*").
		Joins("INNER JOIN sample_metadata ON samples.id = sample_metadata.sample_id").
		Where("parent_sample_id = ?", parentSampleId).
		Find(&samples).Error; err != nil {
		return nil, commonUtils.HandleORMError(err)
	}

	return samples, nil
}
