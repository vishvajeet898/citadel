package constants

// Task Statuses
const (
	TASK_STATUS_PENDING           = "pending"
	TASK_STATUS_IN_PROGRESS       = "in_progress"
	TASK_STATUS_WITHHELD_APPROVAL = "withheld_approval"
	TASK_STATUS_CO_AUTHORIZE      = "co_authorize"
	TASK_STATUS_COMPLETED         = "completed"
)

var TASK_STATUSES = []string{
	TASK_STATUS_PENDING,
	TASK_STATUS_IN_PROGRESS,
	TASK_STATUS_WITHHELD_APPROVAL,
	TASK_STATUS_CO_AUTHORIZE,
	TASK_STATUS_COMPLETED,
}

// Order Types
const (
	ORDER_TYPE_AT_HOME     = "at_home"
	ORDER_TYPE_AT_CLINIC   = "at_clinic"
	ORDER_TYPE_AT_CAMP     = "at_camp"
	ORDER_TYPE_LAB_DROPOFF = "lab_dropoff"
)

var ORDER_TYPES = []string{
	ORDER_TYPE_AT_HOME,
	ORDER_TYPE_AT_CLINIC,
	ORDER_TYPE_AT_CAMP,
	ORDER_TYPE_LAB_DROPOFF,
}

// Calling Types
const (
	CALLING_TYPE_CUSTOMER = "customer"
	CALLING_TYPE_DOCTOR   = "doctor"
)

var CALLING_TYPES = []string{
	CALLING_TYPE_CUSTOMER,
	CALLING_TYPE_DOCTOR,
}

// CollectionType to OrderType Mapping
var CollecTypeToOrderTypeMap = map[uint]string{
	CollectionTypeHomeCollection:    ORDER_TYPE_AT_HOME,
	CollectionTypePickUpFromPartner: ORDER_TYPE_AT_CLINIC,
	CollectionTypeCamps:             ORDER_TYPE_AT_CAMP,
	CollectionTypePartnerDropOff:    ORDER_TYPE_LAB_DROPOFF,
}
