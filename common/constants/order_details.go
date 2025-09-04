package constants

const (
	CollectionTypeHomeCollection    = 0
	CollectionTypePartnerDropOff    = 1
	CollectionTypePickUpFromPartner = 2
	CollectionTypeCamps             = 3
)

const (
	SourcePatientPublicForm = "Patient Public Form"
	SourcePatientApp        = "Patient App"
	SourceDoctorApp         = "Doctor App"
	SourceOMS               = "OMS"
	SourceUnbounce          = "Unbounce"
	SourceTrfApp            = "Trf App"
	SourcePartnerAPI        = "partner_api"
	SourceFrontDeskForm     = "Clinic Front Desk"
	RequestSourcePHC        = "phc"
	SourceDoctorUTMForm     = "Doctor UTM form"
	SourcePartnerUTMForm    = "Partnership UTM form"
	SourcePaperPlane        = "Partner API - Paperplane"
)

const (
	TrfPending    = 1
	TrfTestsAdded = 2
)

const (
	OrderRequestedUint   uint = 1
	OrderOrderedUint     uint = 2
	OrderDoneNotSentUint uint = 3
	OrderCancelledUint   uint = 5
	OrderCompletedUint   uint = 6
)

const (
	OrderRequested   = "requested"
	OrderOrdered     = "ordered"
	OrderDoneNotSent = "done_not_sent"
	OrderCancelled   = "cancelled"
	OrderCompleted   = "completed"
)

var OrderStatusMapUint = map[uint]string{
	OrderRequestedUint:   OrderRequested,
	OrderOrderedUint:     OrderOrdered,
	OrderDoneNotSentUint: OrderDoneNotSent,
	OrderCancelledUint:   OrderCancelled,
	OrderCompletedUint:   OrderCompleted,
}
