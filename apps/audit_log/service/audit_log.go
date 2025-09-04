package service

import (
	mapper "github.com/Orange-Health/citadel/apps/audit_log/mapper"
	"github.com/Orange-Health/citadel/apps/audit_log/structures"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
)

type AuditLogServiceInterface interface {
	GetLogsByOrderId(omsOrderId string) (map[uint]structures.SampleLogBody, *commonStructures.CommonError)
}

func (auditLogService *AuditLogService) GetLogsByOrderId(omsOrderId string) (map[uint]structures.SampleLogBody, *commonStructures.CommonError) {
	sampleAuditLogs, cErr := auditLogService.AuditLogDao.GetSampleAuditLogs(omsOrderId)
	if cErr != nil {
		return nil, cErr
	}
	sampleMetadataAuditLogs, cErr := auditLogService.AuditLogDao.GetSampleMetadataAuditLogs(omsOrderId)
	if cErr != nil {
		return nil, cErr
	}
	userIdToUserNameMap, _ := auditLogService.UserService.GetUserIdNameMap()
	sampleAuditDbLogs, sampleVialTypeMap, sampleParentMap := mapper.MapSampleAuditLogs(sampleAuditLogs, userIdToUserNameMap)
	sampleAuditMetaDbLogs := mapper.MapSampleMetadataAuditLogs(sampleMetadataAuditLogs, userIdToUserNameMap)
	mergedLogs := mapper.MergeAndGroupSampleAuditDbLogsByTimestamp(sampleAuditDbLogs, sampleAuditMetaDbLogs, sampleVialTypeMap, sampleParentMap)
	return mergedLogs, nil
}
