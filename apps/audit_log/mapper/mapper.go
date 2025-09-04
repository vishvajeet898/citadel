package mapper

import (
	"slices"
	"time"

	"github.com/Orange-Health/citadel/apps/audit_log/constants"
	"github.com/Orange-Health/citadel/apps/audit_log/structures"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

func MapSampleAuditLogs(sampleAuditLogs []commonModels.SamplesAudit, usernameMap map[uint]string) (
	map[uint][]structures.SampleDBLog, map[uint]uint, map[uint]uint) {
	// Sort the sample audit logs by log timestamp in descending order
	// This is important to ensure that the latest log is always at the top
	// and we can compare it with the previous log to create the DB logs
	slices.SortFunc(sampleAuditLogs, func(a, b commonModels.SamplesAudit) int {
		if a.LogTimestamp.After(b.LogTimestamp) {
			return -1
		} else if a.LogTimestamp.Before(b.LogTimestamp) {
			return 1
		}
		return 0
	})
	sampleAuditLogsMap := make(map[uint][]commonModels.SamplesAudit)
	sampleDbLogsMap := make(map[uint][]structures.SampleDBLog)
	sampleVialTypeMap := make(map[uint]uint)
	parentSampleMap := make(map[uint]uint)
	for _, sampleAudit := range sampleAuditLogs {
		if sampleAudit.ParentSampleId != 0 {
			parentSampleMap[sampleAudit.Id] = sampleAudit.ParentSampleId
		}
		sampleAuditLogsMap[sampleAudit.Id] = append(sampleAuditLogsMap[sampleAudit.Id], sampleAudit)
		if sampleAudit.VialTypeId != 0 {
			sampleVialTypeMap[sampleAudit.Id] = sampleAudit.VialTypeId
		}
	}
	for sampleId, auditLogs := range sampleAuditLogsMap {
		dbLogs := []structures.SampleDBLog{}
		for index := range auditLogs {
			if index < len(auditLogs)-1 {
				compareAndCreateDbLogs(auditLogs[index], auditLogs[index+1], &dbLogs, usernameMap)
			} else if _, exists := parentSampleMap[sampleId]; !exists {
				// If the sample is not a parent sample, we need to add the last log
				// as a DB log. This is because the last log is not compared with
				// any other log and we need to add it to the DB logs.
				dbLogs = append(dbLogs, getDbLogStruct(auditLogs[index].LogAction, "", "", "", auditLogs[index].LogTimestamp,
					auditLogs[index].UpdatedBy, usernameMap))
			}
		}
		sampleDbLogsMap[sampleId] = dbLogs
	}
	return sampleDbLogsMap, sampleVialTypeMap, parentSampleMap
}

func getUserName(userId uint, userIdToUserNameMap map[uint]string) string {
	if userId == commonConstants.CitadelSystemId {
		return commonConstants.SuperlabSystemName
	}
	if userName, exists := userIdToUserNameMap[userId]; exists {
		return userName
	}
	return ""
}

func compareAndCreateDbLogs(newSampleAudit, oldSampleAudit commonModels.SamplesAudit, dbLogs *[]structures.SampleDBLog,
	usernameMap map[uint]string) {
	operation := newSampleAudit.LogAction
	if newSampleAudit.VisitId != oldSampleAudit.VisitId {
		*dbLogs = append(*dbLogs, getDbLogStruct(operation, constants.LogFieldVisitID, oldSampleAudit.VisitId,
			newSampleAudit.VisitId, newSampleAudit.LogTimestamp, newSampleAudit.UpdatedBy, usernameMap))
	}

	if newSampleAudit.Barcode != oldSampleAudit.Barcode {
		*dbLogs = append(*dbLogs, getDbLogStruct(operation, constants.LogFieldBarcode, oldSampleAudit.Barcode,
			newSampleAudit.Barcode, newSampleAudit.LogTimestamp, newSampleAudit.UpdatedBy, usernameMap))
	}

	if newSampleAudit.Status != oldSampleAudit.Status {
		*dbLogs = append(*dbLogs, getDbLogStruct(operation, constants.LogFieldStatus, oldSampleAudit.Status,
			newSampleAudit.Status, newSampleAudit.LogTimestamp, newSampleAudit.UpdatedBy, usernameMap))
	}

	if newSampleAudit.LabId != oldSampleAudit.LabId {
		*dbLogs = append(*dbLogs, getDbLogStruct(operation, constants.LogFieldLabID,
			commonUtils.ConvertUintToString(oldSampleAudit.LabId), commonUtils.ConvertUintToString(newSampleAudit.LabId),
			newSampleAudit.LogTimestamp, newSampleAudit.UpdatedBy, usernameMap))
	}

	if newSampleAudit.DestinationLabId != oldSampleAudit.DestinationLabId {
		*dbLogs = append(*dbLogs, getDbLogStruct(operation, constants.LogFieldDestinationLabID,
			commonUtils.ConvertUintToString(oldSampleAudit.DestinationLabId),
			commonUtils.ConvertUintToString(newSampleAudit.DestinationLabId), newSampleAudit.LogTimestamp,
			newSampleAudit.UpdatedBy, usernameMap))
	}

	if newSampleAudit.RejectionReason != oldSampleAudit.RejectionReason {
		*dbLogs = append(*dbLogs, getDbLogStruct(operation, constants.LogFieldRejectionReason,
			oldSampleAudit.RejectionReason, newSampleAudit.RejectionReason, newSampleAudit.LogTimestamp,
			newSampleAudit.UpdatedBy, usernameMap))
	}

	if newSampleAudit.ParentSampleId != oldSampleAudit.ParentSampleId {
		*dbLogs = append(*dbLogs, getDbLogStruct(operation, constants.LogFieldParentSampleID,
			commonUtils.ConvertUintToString(oldSampleAudit.ParentSampleId),
			commonUtils.ConvertUintToString(newSampleAudit.ParentSampleId), newSampleAudit.LogTimestamp,
			newSampleAudit.UpdatedBy, usernameMap))
	}

	if newSampleAudit.DeletedAt != oldSampleAudit.DeletedAt {
		*dbLogs = append(*dbLogs, getDbLogStruct(operation, constants.LogFieldDeletedAt,
			commonUtils.GetDeletedAtString(*oldSampleAudit.DeletedAt),
			commonUtils.GetDeletedAtString(*newSampleAudit.DeletedAt), newSampleAudit.LogTimestamp,
			newSampleAudit.UpdatedBy, usernameMap))
	}
}

func MapSampleMetadataAuditLogs(sampleMetadataAuditLogs []commonModels.SampleMetadataAudit,
	usernameMap map[uint]string) map[uint][]structures.SampleDBLog {
	slices.SortFunc(sampleMetadataAuditLogs, func(a, b commonModels.SampleMetadataAudit) int {
		if a.LogTimestamp.After(b.LogTimestamp) {
			return -1
		} else if a.LogTimestamp.Before(b.LogTimestamp) {
			return 1
		}
		return 0
	})
	sampleMetadataAuditLogsMap := make(map[uint][]commonModels.SampleMetadataAudit)
	sampleDbLogsMap := make(map[uint][]structures.SampleDBLog)
	for _, sampleMetadataAudit := range sampleMetadataAuditLogs {
		sampleMetadataAuditLogsMap[sampleMetadataAudit.SampleId] = append(
			sampleMetadataAuditLogsMap[sampleMetadataAudit.SampleId], sampleMetadataAudit)
	}
	for sampleId, auditLogs := range sampleMetadataAuditLogsMap {
		dbLogs := []structures.SampleDBLog{}
		for index := range auditLogs {
			if index < len(auditLogs)-1 {
				compareAndCreateSampleMetaDataDbLogs(auditLogs[index], auditLogs[index+1], &dbLogs, usernameMap)
			}
		}
		sampleDbLogsMap[sampleId] = dbLogs
	}
	return sampleDbLogsMap
}

func compareAndCreateSampleMetaDataDbLogs(newSampleMetadataAudit, oldSampleMetadataAudit commonModels.SampleMetadataAudit,
	dbLogs *[]structures.SampleDBLog, usernameMap map[uint]string) {
	operation := newSampleMetadataAudit.LogAction
	if newSampleMetadataAudit.NotReceivedReason != oldSampleMetadataAudit.NotReceivedReason {
		*dbLogs = append(*dbLogs, getDbLogStruct(operation, constants.LogFieldNotReceivedReason,
			oldSampleMetadataAudit.NotReceivedReason, newSampleMetadataAudit.NotReceivedReason,
			newSampleMetadataAudit.LogTimestamp, newSampleMetadataAudit.UpdatedBy, usernameMap))
	}

	if newSampleMetadataAudit.CollectLaterReason != oldSampleMetadataAudit.CollectLaterReason {
		*dbLogs = append(*dbLogs, getDbLogStruct(operation, constants.LogFieldCollectLaterReason,
			oldSampleMetadataAudit.CollectLaterReason, newSampleMetadataAudit.CollectLaterReason,
			newSampleMetadataAudit.LogTimestamp, newSampleMetadataAudit.UpdatedBy, usernameMap))
	}

	if newSampleMetadataAudit.RejectingLab != oldSampleMetadataAudit.RejectingLab {
		*dbLogs = append(*dbLogs, getDbLogStruct(operation, constants.LogFieldRejectingLab,
			commonUtils.ConvertUintToString(oldSampleMetadataAudit.RejectingLab),
			commonUtils.ConvertUintToString(newSampleMetadataAudit.RejectingLab), newSampleMetadataAudit.LogTimestamp,
			newSampleMetadataAudit.UpdatedBy, usernameMap))
	}
}

func getDbLogStruct(operation, fieldName, oldValue, newValue string, logTimestamp time.Time,
	userId uint, userNameMap map[uint]string) structures.SampleDBLog {
	userName := getUserName(userId, userNameMap)
	return structures.SampleDBLog{
		Operation:    operation,
		FieldName:    fieldName,
		OldValue:     oldValue,
		NewValue:     newValue,
		UserName:     userName,
		LogTimestamp: logTimestamp,
	}
}

func MergeAndGroupSampleAuditDbLogsByTimestamp(sampleLogs, sampleMetaLogs map[uint][]structures.SampleDBLog,
	sampleVialMap, sampleParentMap map[uint]uint) map[uint]structures.SampleLogBody {
	mergedLogs := make(map[uint]map[time.Time][]structures.SampleDBLog)
	for sampleId, sampleDbLogs := range sampleLogs {
		mergedLogs[sampleId] = make(map[time.Time][]structures.SampleDBLog)
		for _, dbLog := range sampleDbLogs {
			mergedLogs[sampleId][dbLog.LogTimestamp] = append(mergedLogs[sampleId][dbLog.LogTimestamp], dbLog)
		}
	}
	for sampleId, sampleMetaDbLogs := range sampleMetaLogs {
		if _, exists := mergedLogs[sampleId]; !exists {
			mergedLogs[sampleId] = make(map[time.Time][]structures.SampleDBLog)
		}
		for _, dbLog := range sampleMetaDbLogs {
			mergedLogs[sampleId][dbLog.LogTimestamp] = append(mergedLogs[sampleId][dbLog.LogTimestamp], dbLog)
		}
	}
	for sampleId, dbLogs := range mergedLogs {
		if _, exists := sampleParentMap[sampleId]; exists {
			for timestamp, logs := range dbLogs {
				mergedLogs[sampleParentMap[sampleId]][timestamp] = append(mergedLogs[sampleParentMap[sampleId]][timestamp], logs...)
			}
			delete(mergedLogs, sampleId)
		}
	}
	sampleLogBody := make(map[uint]structures.SampleLogBody)
	for sampleId, dbLogs := range mergedLogs {
		sampleLogBody[sampleId] = structures.SampleLogBody{
			SampleId:   sampleId,
			VialTypeId: sampleVialMap[sampleId],
			Logs:       dbLogs,
		}
	}
	return sampleLogBody

}
