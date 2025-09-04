package dao

import (
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

type DataLayer interface {
	GetSampleAuditLogs(omsOrderId string) ([]commonModels.SamplesAudit, *commonStructures.CommonError)
	GetSampleMetadataAuditLogs(omsOrderId string) ([]commonModels.SampleMetadataAudit, *commonStructures.CommonError)
}

func (auditLogDao *AuditLogDao) GetSampleAuditLogs(omsOrderId string) ([]commonModels.SamplesAudit, *commonStructures.CommonError) {
	sampleAudits := []commonModels.SamplesAudit{}
	if err := auditLogDao.Db.Where("oms_order_id = ?", omsOrderId).Find(&sampleAudits).Error; err != nil {
		return sampleAudits, commonUtils.HandleORMError(err)
	}
	return sampleAudits, nil
}

func (auditLogDao *AuditLogDao) GetSampleMetadataAuditLogs(omsOrderId string) ([]commonModels.SampleMetadataAudit, *commonStructures.CommonError) {
	sampleMetadataAudits := []commonModels.SampleMetadataAudit{}
	if err := auditLogDao.Db.Where("oms_order_id = ?", omsOrderId).Find(&sampleMetadataAudits).Error; err != nil {
		return sampleMetadataAudits, commonUtils.HandleORMError(err)
	}
	return sampleMetadataAudits, nil
}
