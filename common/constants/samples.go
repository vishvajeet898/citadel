package constants

const (
	SampleNotReceived        = "not_received"
	SampleRejected           = "rejected"
	SamplePartiallyRejected  = "partially_rejected"
	SampleSynced             = "synced"
	SampleReceived           = "received"
	SampleDeleted            = "deleted"
	SampleCollectionDone     = "collection_done"
	SampleDefault            = "default"
	SampleAccessioned        = "accessioned"
	SampleNotCollectedEmedic = "not_collected_emedic"
	SampleTransferred        = "transferred"
	SampleTransferFailed     = "transfer_failed"
	SampleInTransfer         = "in_transfer"
	SampleOutsourced         = "outsourced"
)

const (
	SampleNotReceivedUint        uint = 1
	SampleRejectedUint           uint = 2
	SamplePartiallyRejectedUint  uint = 3
	SampleSyncedUint             uint = 4
	SampleReceivedUint           uint = 5
	SampleDeletedUint            uint = 6
	SampleCollectionDoneUint     uint = 7
	SampleDefaultUint            uint = 8
	SampleAccessionedUint        uint = 9
	SampleNotCollectedEmedicUint uint = 10
	SampleTransferredUint        uint = 11
	SampleTransferFailedUint     uint = 12
	SampleOutsourcedUint         uint = 13
	SampleInTransferUint         uint = 14
)

var SampleStatusMap = map[string]uint{
	SampleNotReceived:        SampleNotReceivedUint,
	SampleRejected:           SampleRejectedUint,
	SamplePartiallyRejected:  SamplePartiallyRejectedUint,
	SampleSynced:             SampleSyncedUint,
	SampleReceived:           SampleReceivedUint,
	SampleDeleted:            SampleDeletedUint,
	SampleCollectionDone:     SampleCollectionDoneUint,
	SampleDefault:            SampleDefaultUint,
	SampleAccessioned:        SampleAccessionedUint,
	SampleNotCollectedEmedic: SampleNotCollectedEmedicUint,
	SampleTransferred:        SampleTransferredUint,
	SampleTransferFailed:     SampleTransferFailedUint,
	SampleOutsourced:         SampleOutsourcedUint,
	SampleInTransfer:         SampleInTransferUint,
}

var SampleStatusMapReverse = map[uint]string{
	SampleNotReceivedUint:        SampleNotReceived,
	SampleRejectedUint:           SampleRejected,
	SamplePartiallyRejectedUint:  SamplePartiallyRejected,
	SampleSyncedUint:             SampleSynced,
	SampleReceivedUint:           SampleReceived,
	SampleDeletedUint:            SampleDeleted,
	SampleCollectionDoneUint:     SampleCollectionDone,
	SampleDefaultUint:            SampleDefault,
	SampleAccessionedUint:        SampleAccessioned,
	SampleNotCollectedEmedicUint: SampleNotCollectedEmedic,
	SampleTransferredUint:        SampleTransferred,
	SampleTransferFailedUint:     SampleTransferFailed,
	SampleOutsourcedUint:         SampleOutsourced,
	SampleInTransferUint:         SampleOutsourced,
}

var (
	MultipleOrdersReceivingEnabled = Config.GetBool("receiving_desk.multiple_orders_enabled")
)

var (
	QnsEnabledStates        = []string{SampleReceived, SamplePartiallyRejected, SampleSynced, SampleAccessioned}
	LisFilteredSampleStatus = []string{SampleNotReceived, SampleRejected, SampleSynced, SampleAccessioned, SampleNotCollectedEmedic}
)

var (
	NotReceivedReasonCollectionRefused = "COLLECTION_REFUSED"
	NotReceivedReasonMissingSample     = "MISSING_SAMPLE"
	NotReceivedReasonOthers            = "OTHER"
	NotCollectedReasonCollectLater     = "Sample to be collected later"
	PatientDeniedSampleCollection      = "Patient denied giving sample"
	AutoCancellationDueToRNRReason     = "AUTO_CANCELLATION_RNR"
)

var (
	OmsTaskTypePrimaryCollection                uint = 1
	OmsTaskTypeRecollection                     uint = 2
	OmsTaskTypeDigitalTRFCollection             uint = 4
	OmsTaskTypeDigitalTRFRecollectionFromRunner uint = 5
)

const (
	AddNewSampleDefaultType   = "add_new_sample_default"
	AddNewSampleCollectedType = "add_new_sample_collected"
	MapExistingSampleType     = "map_existing_sample"
	LabIdModificationType     = "lab_id_modification"
)

var (
	VialTypesToBeSkippedForReverseLogistics = []uint{40, 43, 44, 45, 46, 47}
)
