package constants

// Task types
const (
	TASK_TYPE_CRITICAL     = "critical"
	TASK_TYPE_WITHHELD     = "withheld"
	TASK_TYPE_CO_AUTHORIZE = "co_authorize"
	TASK_TYPE_NORMAL       = "normal"
	TASK_TYPE_IN_PROGRESS  = "in_progress"
	TASK_TYPE_AMENDMENT    = "amendment"
)

var TASK_TYPES = []string{
	TASK_TYPE_CRITICAL,
	TASK_TYPE_WITHHELD,
	TASK_TYPE_CO_AUTHORIZE,
	TASK_TYPE_NORMAL,
	TASK_TYPE_IN_PROGRESS,
	TASK_TYPE_AMENDMENT,
}

const (
	SPECIAL_REQUIREMENT_MORPHLE = "contains_morphle"
	SPECIAL_REQUIREMENT_PACKAGE = "contains_package"
)

var SPECIAL_REQUIREMENTS = []string{
	SPECIAL_REQUIREMENT_MORPHLE,
	SPECIAL_REQUIREMENT_PACKAGE,
}

const (
	ERROR_WHILE_FETCHING_TEST_DETAILS  = "error while fetching test details"
	ERROR_WHILE_FETCHING_ORDER_DETAILS = "error while fetching order details"
	ERROR_WHILE_FETCHING_VISIT_DETAILS = "error while fetching visit details"
	ERROR_NO_TEST_FOUND_IN_BARCODE     = "no test found in barcode"
	ERROR_NO_INHOUSE_TESTS_FOUND       = "no inhouse tests found in this sample"
)

const (
	ServiceTypeScan    = "scan"
	ServiceTypeDetails = "details"
)

const (
	ERROR_INVALID_ORDER_ID = "invalid order id"
	ERROR_INVALID_BARCODE  = "invalid barcode"
	ERROR_INVALID_VISIT_ID = "invalid visit id"
)

const (
	INFO_SCREEN_SEARCH_TYPE_ORDER_ID = "order_id"
	INFO_SCREEN_SEARCH_TYPE_VISIT_ID = "visit_id"
	INFO_SCREEN_SEARCH_TYPE_BARCODE  = "barcode"
)
