package constants

// Common Error Messages
const (
	ERROR_FAILED_TO_UNMARSHAL_JSON       = "failed to unmarshal json"
	ERROR_INVALID_LIMIT                  = "invalid limit provided"
	ERROR_FAILED_TO_MARSHAL_PAYLOAD      = "failed to marshal payload"
	ERROR_FAILED_TO_DECODE_BASE64_STRING = "failed to decode base64 string"
	TECHNICAL_ERROR                      = "Technical error! Please try again"
	ERROR_FAILED_TO_GET_LOCAL_LOCATION   = "failed to get local location"
	ERROR_WHILE_PARSING_FLOAT_VALUE      = "Error while parsing float value"
	INVALID_INVESTIGATION_ENTERED_TIME   = "Invalid investigation entered at time"
)

// Worker Error messages
const (
	ERROR_WHILE_FETCHING_PAYLOAD_FROM_REDIS        = "error while fetching payload from redis"
	ERROR_WHILE_PARSING_EVENT_PAYLOAD              = "error while parsing event payload"
	ERROR_FAILED_TO_PROCESS_COMPLETED_EVENT        = "failed to process completed event"
	ERROR_FAILED_TO_PROCESS_APPROVE_EVENT          = "failed to process approve event"
	ERROR_FAILED_TO_GET_MESSAGE_FROM_EVENT_PAYLOAD = "failed to get message from event payload"
	ERROR_MANUAL_REPORT_UPLOAD_TASK_IN_PROGRESS    = "manual report upload task in progress"
	ERROR_LIS_EVENT_TASK_IN_PROGRESS               = "lis event task in progress"
	ERROR_COMPLETED_EVENT_TASK_IN_PROGRESS         = "completed event task in progress"
	ERROR_APPROVE_EVENT_TASK_IN_PROGRESS           = "approve event task in progress"
	ERROR_ATTACHMENT_EVENT_TASK_IN_PROGRESS        = "attachment event task in progress"
	ERROR_RERUN_EVENT_TASK_IN_PROGRESS             = "rerun event task in progress"
	ERROR_CREATE_UPDATE_ORDER_TASK_IN_PROGRESS     = "create update order task in progress"
	ERROR_COMPLETED_ORDER_TASK_IN_PROGRESS         = "completed order task in progress"
	ERROR_CANCELLED_ORDER_TASK_IN_PROGRESS         = "cancelled order task in progress"
	ERROR_FAILED_TO_START_SERVER                   = "failed to start server"
	ERROR_FAILED_TO_SEND_TASK                      = "failed to send task"
	ERROR_REPORT_PDF_FORMAT_IGNORED                = "report pdf format ignored"
	ERROR_EMPTY_MESSAGE_IN_EVENT_PAYLOAD           = "empty message in event payload"
)

// Redis Error messages
const (
	ERROR_KEY_DOES_NOT_EXIST               = "key does not exist"
	ERROR_WHILE_FETCHING_DATA_FROM_REDIS   = "error while fetching data from redis"
	ERROR_FAILED_TO_SET_REPORT_RELEASE_KEY = "failed to set report release key"
	ERROR_FAILED_TO_SET_TASK_PATH_MAPPING  = "failed to set task pathologist mapping"
)

// DB Error messages
const (
	ERROR_WHILE_CONNECTING_TO_DATABASE = "error while connecting to database"
	ERROR_IN_TRANSACTION               = "error in transaction"
	ERROR_NO_RECORDS_FOUND             = "%s : No records found"
)

// Sentry Error Messages
const (
	SENTRY_INITIALIZATION_FAILED = "sentry initialization failed"
)

// SQS Error Messages
const (
	ERROR_CREATING_QUEUE = "error creating queue"
)

// API Client Error Messages
const (
	HTTP_REQUEST_FAILED = "http request failed"
)

// AWS Adapter Error Messages
const (
	ERROR_WHILE_CREATING_AWS_PUB_SUB_ADAPTER = "error while creating aws pub sub adapter"
	ERROR_WHILE_POLLING_MESSAGES             = "error while polling messages"
	ERROR_WHILE_PUBLISHING_MESSAGE           = "error while publishing message on AWS"
)

// Attune Error Messages
const (
	ERROR_WHILE_LOGGING_INTO_ATTUNE            = "error while logging into attune"
	ERROR_WHILE_GETTING_ORDER_DATA_FROM_ATTUNE = "error while getting order data from attune"
	ERROR_WHILE_SENDING_ATTUNE_WEBHOOK         = "error while sending attune webhook"
	ERROR_LAB_VISIT_ID_NOT_FOUND_IN_ATTUNE     = "lab visit id not found in attune"
	ERROR_GETTING_ATTUNE_ORG_CODE              = "error while getting attune org code"
	ERROR_NO_TESTS_TO_BE_SYNCED                = "no tests to be synced to lis"
	ERROR_FAILED_TO_GET_VISIT_ID_FROM_REDIS    = "failed to get visit count from redis"
	ERROR_FAILED_TO_SET_VISIT_ID_IN_REDIS      = "failed to set visit count in redis"
	ERROR_LIS_SYNC_ALREADY_IN_PROGRESS         = "lis sync already in progress"
	ERROR_LIS_ORDER_ID_NOT_FOUND               = "lis order id not found"
	ERROR_LIS_ORDER_INFO_NOT_FOUND             = "lis order info not found"
	ERROR_INVALID_REPORT_PDF_FORMAT            = "invalid report pdf format"
	ERROR_REPORT_PDF_NOT_FOUND                 = "report pdf not found"
)

// Tasks Error Messages
const (
	ERROR_WHILE_FETCHING_TASKS               = "error while fetching tasks"
	ERROR_INVALID_TASK_ID                    = "invalid task id"
	ERROR_INVALID_VISIT_ID                   = "invalid visit id"
	ERROR_INVALID_TASK_TYPE                  = "invalid task type"
	ERROR_INVALID_TASK_STATUS                = "invalid task status"
	ERROR_INVALID_SPECIAL_REQUIREMENT        = "invalid special requirement"
	ERROR_INVALID_ORDER_TYPE                 = "invalid order type"
	ERROR_INVALID_LAB_ID                     = "invalid lab id"
	ERROR_INVALID_REQUEST_ID                 = "invalid request id"
	ERROR_INVALID_ORDER_ID                   = "invalid order id"
	ERROR_IN_UNDOING_REPORT_RELEASE          = "error in undoing report release"
	ERROR_ACTION_NOT_ALLOWED_FOR_PATHOLOGIST = "action not allowed for pathologist"
	ERROR_INVALID_CALLING_TYPE               = "invalid calling type"
	ERROR_CO_AUTHORIZE_TO_MISMATCH           = "can't coauthorize two different doctors"
	ERROR_CO_AUTHORIZE_TO_MISSING            = "coauthorize to is missing"
	ERROR_CO_AUTHORIZE_TO_SELF               = "can't coauthorize to self"
	ERROR_TASK_NOT_FOUND                     = "task not found"
	ERROR_FAILED_TO_CREATE_UPDATE_TASK       = "failed to create or update task"
)

// Task Pathologist Mapping Error Messages
const (
	ERROR_TPM_EXISTS        = "Task Pathologist Mapping already exists for the given Task ID"
	ERROR_TPM_IN_SAME_STATE = "Task Pathologist Mapping already exists in %s state"
)

// Test Details Error Messages
const (
	ERROR_INVALID_TEST_DETAILS_ID                        = "invalid test details id"
	ERROR_INVALID_TEST_STATUS                            = "invalid test status"
	ERROR_TEST_STATUS_WITHHELD_CO_AUTHORIZE_NON_COUPLING = "test status withheld can not be coupled with co authorize"
	ERROR_NO_TESTS_IDS_TO_DELETE                         = "no test ids to delete"
	ERROR_TEST_NOT_FOUND                                 = "test not found"
	ERROR_NO_TEST_DETAILS_FOUND                          = "no test details found"
	ERROR_INVALID_PROCESSING_LAB_REQUEST                 = "invalid processing lab request"
	ERROR_PROCESSING_LAB_IS_INHOUSE                      = "new processing lab should be outsourced"
)

// Patient Error Messages
const (
	ERROR_INVALID_PATIENT_ID = "invalid patient id"
)

// Investigations Error Messages
const (
	ERROR_INVALID_MASTER_INVESTIGATION_IDS    = "invalid master investigation ids"
	ERROR_INVALID_INVESTIGATION_CODE          = "invalid investigation code"
	ERROR_INVALID_INVESTIGATION_VALUE         = "invalid investigation value"
	ERROR_NEGATIVE_INVESTIGATION_VALUE        = "Negative Value Present in investigation"
	ERROR_INVALID_MASTER_INVESTIGATION_ID     = "invalid master investigation id"
	ERROR_INVALID_INVESTIGATION_RESULT_ID     = "invalid investigation result id"
	ERROR_INVALID_INVESTIGATION_STATUS        = "invalid investigation status"
	ERROR_INVESTIGATION_IDS_MISMATCH          = "investigation ids mismatch"
	ERROR_INVESTIGATION_STATUS_MISMATCH       = "investigation status mismatch"
	ERROR_RERUN_STATUS_COUPLING               = "rerun status can only be coupled with approve status"
	ERROR_NO_INVESTIGATION_RESULTS_FOUND      = "no investigation results found"
	ERROR_NO_FORMULA_FOUND                    = "Calculation formula not found"
	ERROR_INVALID_RESULT                      = "invalid result"
	ERROR_IN_CALCULATING_RESULT               = "Error occurred while calculating result. Please try again"
	ERROR_IN_GETTING_LAST_INVESTIGATION_VALUE = "Error in getting last investigation value"
	ERROR_IN_GETTING_QC_FAILED_RERUN_DATA     = "Error in getting QC failed rerun data"
)

// Templates Error Messages
const (
	ERROR_INVALID_TEMPLATE_TYPE = "invalid template type"
)

// User Error Messages
const (
	INVALID_USER_ID_RECEIVED = "invalid user id received"
	USER_ID_NOT_FOUND_IN_JWT = "user id not found in jwt"
	INVALID_USER_TYPE        = "invalid user type"
	USER_NAME_NOT_FOUND      = "user name not found"
	EMAIL_NOT_FOUND          = "email not found"
	USER_ALREADY_EXISTS      = "user already exists"
)

// Oms Error Messages
const (
	ERROR_WHILE_FETCHING_REPORT_GENERATION_EVENT_DETAILS = "error in fetching report generation event details from OMS"
	ERROR_NO_REPORT_GENERATION_EVENT_FOUND               = "no report generation event found"
	ERROR_WHILE_UPDATING_PATIENT_DETAILS                 = "error while updating patient details"
	ERROR_WHILE_FETCHING_VISIT_DETAILS_BY_TEST_IDS       = "error while fetching visit details by test ids"
	ERROR_WHILE_FETCHING_DELTA_VALUES                    = "error while fetching delta values"
	ERROR_WHILE_FETCHING_PATIENT_PAST_RECORDS            = "error while fetching patient past records"
	ERROR_WHILE_FETCHING_ORDER_DETAILS                   = "error while fetching order details"
)

// Cds Error Messages
const (
	ERROR_IN_GETTING_PANEL_DETAILS               = "error in getting panel details"
	ERROR_IN_GETTING_INVESTIGATION_DETAILS       = "error in getting investigation details"
	ERROR_IN_GETTING_DEPARTMENT_MAPPING          = "error in getting department mapping"
	ERROR_IN_GETTING_DEPENDENT_INVESTIGATIONS    = "error in getting dependent investigations"
	ERROR_IN_GETTING_VIALS                       = "error in getting vials"
	ERROR_IN_GETTING_LABS                        = "error in getting labs"
	ERROR_IN_GETTING_COLLECTION_SEQUENCE         = "error in getting collection sequence"
	ERROR_IN_GETTING_MASTER_TESTS                = "error in getting master tests"
	ERROR_MASTER_TEST_NOT_FOUND                  = "master test not found"
	ERROR_WHILE_GETTING_NON_INHOUSE_REPORT_TESTS = "error while getting non inhouse tests"
	ERROR_TESTS_AND_PACKAGES_CANNOT_BE_EMPTY     = "tests and packages cannot be empty"
	ERROR_IN_GETTING_DEDUPLICATION_RESPONSE      = "error in getting deduplication response"
	ERROR_IN_GETTING_NRL_ENABLED_TESTS           = "error in getting nrl enabled tests"
)

// Partner Api Error Messages
const (
	ERROR_WHILE_GETTING_PARTNER_DETAILS      = "error while getting partner details"
	ERROR_WHILE_GETTING_BULK_PARTNER_DETAILS = "error while getting bulk partner details"
	ERROR_PARTNER_DOES_NOT_EXIST             = "partner does not exist"
)

// Health Api Error Messages
const (
	ERROR_WHILE_GETTING_DOCTORS               = "error while getting doctors"
	ERROR_DOCTOR_DOES_NOT_EXIST               = "doctor does not exist"
	ERROR_WHILE_SENDING_GENERIC_SLACK_MESSAGE = "error while sending generic slack message"
)

// Accounts Api Error Messages
const (
	ERROR_WHILE_CREATING_FRESHDESK_TICKET = "error while creating freshdesk ticket"
)

// Report Rebranding Error Messages
const (
	ERROR_WHILE_RESIZING_MEDIA            = "error while resizing media"
	ERROR_WHILE_ATTACHING_COBRANDED_IMAGE = "error while attaching cobranded image"
)

// Attachments Error Messages
const (
	ATTACHMENT_TYPE_REQUIRED     = "attachment type is required"
	ATTACHMENT_LABEL_REQUIRED    = "attachment label is required"
	FILE_NAME_REQUIRED           = "file name is required"
	INVALID_ATTACHMENT_LABEL     = "invalid attachment label"
	INVALID_ATTACHMENT_TYPE      = "invalid attachment type"
	INVALID_FILE_NAME            = "invalid file name"
	INVALID_ATTACHMENT_EXTENSION = "invalid attachment extension"
	TOKENIZED_URL_ERROR          = "error while getting tokenized file url"
	ERROR_FAILED_TO_DELETE_FILE  = "failed to delete file"
)

// Patient Service Error Messages
const (
	ERROR_WHILE_GETTING_SIMILAR_PATIENT_DETAILS = "error while getting similar patient details"
	ERROR_WHILE_GETTING_PATIENT_DETAILS         = "error while getting patient details"
)

// Samples Error Messages
const (
	ERROR_BARCODE_EXISTS                     = "barcode %s already exists, same barcode can be used within an order for same vial types"
	ERROR_BARCODE_EXISTS_IN_SYSTEM           = "barcode %s already exists in the system"
	ERROR_DUPLICATE_BARCODE                  = "duplicate barcode found"
	ERROR_SAMPLE_NUMBER_EXISTS               = "sample number already exists for this order"
	ERROR_WHILE_UPDATING_SAMPLE              = "error while updating sample"
	ERROR_WHILE_MERGING_SAMPLE               = "error while merging samples"
	ERROR_SAMPLE_ID_CANNOT_BE_ZERO           = "sample id cannot be 0"
	ERROR_SAMPLE_NOT_FOUND                   = "sample/s not found"
	ERROR_BARCODE_NOT_FOUND                  = "barcode not found"
	ERROR_NO_SAMPLES_TO_UPDATE_BARCODE       = "no samples to update barcode"
	ERROR_SAMPLE_INVALID_FOR_REJECTION       = "sample is invalid for rejection"
	ERROR_TEST_IDS_REQUIRED_FOR_RECOLLECTION = "at least 1 test is required for recollection cases"
	ERROR_NO_SAMPLES_FOUND_TO_MARK_COLLECTED = "no samples found to mark collected"
	ERROR_FAILED_TO_DELETE_SAMPLE            = "failed to delete sample"
	ERROR_NO_SAMPLES_FOUND                   = "no samples found"
	ERROR_MULTIPLE_ORDERS_NOT_ALLOWED        = "multiple orders not allowed"
	ERROR_MERGE_SAMPLE_DETAILS               = "error while merging sample details"
	ERROR_SAMPLE_ID_REQUIRED                 = "sample id is required"
	ERROR_NO_INTERLAB_SAMPLES                = "no interlab samples found"
	ERROR_NO_SAMPLE_TEST_MAPPING_DETAILS     = "no sample test mapping details found"
	ERROR_BARCODE_REQUIRED                   = "barcode is required"
	ERROR_VISIT_ID_REQUIRED                  = "visit id is required"
	ERROR_VISIT_ID_NOT_FOUND                 = "visit id not found"
	ERROR_LIS_SYNC_NOT_ENABLED_FOR_LAB       = "lis sync is not enabled for the lab"
	ERROR_SYNCHRONIZING_TASKS                = "error synchronizing tasks"
	ERROR_TASK_SEQUENCE_REQUIRED             = "task sequence is required"
	ERROR_NO_COVID_TEST_VISIT_IDS_FOUND      = "no covid test visit ids found for the given order id"
)

// Vials Error Messages
const (
	ERROR_VIAL_NOT_FOUND = "vial not found"
)

// Labs Error Messages
const (
	ERROR_LAB_NOT_FOUND         = "lab not found"
	LAB_ID_NOT_FOUND_IN_JWT     = "lab id not found in jwt"
	ERROR_LAB_ID_CANNOT_BE_ZERO = "lab id cannot be 0"
	ERROR_LAB_DETAILS_NOT_FOUND = "lab details not found"
)

// Order Details Error Messages
const (
	ERROR_SERVICING_LAB_NOT_FOUND               = "servicing lab not found"
	ERROR_ORDER_ID_NOT_FOUND                    = "order id not found"
	ERROR_ORDER_ID_CANNOT_BE_ZERO               = "order id cannot be 0"
	ERROR_TRF_ID_NOT_FOUND                      = "trf id not found"
	ERROR_INVALID_OMS_ORDER_CREATE_UPDATE_EVENT = "invalid oms order create update event"
)

// Report Generation Error Messages
const (
	ERROR_REPORT_GENRATION_DISABLED    = "report generation is disabled for this city"
	ERROR_NO_VISIT_TESTS_FOUND         = "no visit tests found"
	ERROR_NO_INHOUSE_VISIT_TESTS_FOUND = "no inhouse visit tests found"
)

// Test Sample Mapping Error Messages
const (
	ERROR_TSM_NOT_FOUND                    = "test Sample Mapping not found"
	ERROR_TEST_ID_CANNOT_BE_ZERO           = "test id cannot be 0"
	ERROR_SAMPLE_NUMBER_CANNOT_BE_ZERO     = "sample number cannot be 0"
	ERROR_NO_TEST_SAMPLE_MAPPING_FOUND     = "no test sample mapping found"
	ERROR_NO_TESTS_TO_UPDATE_SAMPLE_NUMBER = "no tests to update sample number"
)

// External Investigations Error Messages
const (
	SYSTEM_EXTERNAL_INVESTIGATION_IDS_REQUIRED = "system_external_investigation_ids query is missing"
	DELETED_BY_QUERY_REQUIRED                  = "deleted_by query is missing"
	EXT_INV_NOT_FOUND_FOR_SYSTEM_EXT_INV_IDS   = "no external investigations found with system_external_investigation_result_id = %v"
	CONTACT_ID_QUERY_REQUIRED                  = "contact_id query is missing"
	LOINC_OR_INV_ID_QUERY_REQUIRED             = "loinc_code and master_investigation_method_mapping_id both are missing. Please provide at least one"
)

// Receiving Desk Error Messages
const (
	ERROR_IN_RECEIVING_FOR_THIS_ORDER_ID   = "error in receiving for this order id"
	ERROR_RECEIVE_AND_SYNC_IN_PROGRESS     = "receive and sync is already in progress for tbis order"
	ERROR_RECEIVING_DESK_DISABLED_FOR_CITY = "receiving desk is disabled for this city"
)

// Contact Error Messages
const ERROR_FAILED_TO_MERGE_INVESTIGATIONS = "Failed to merge external investigations"

const (
	ERROR_FAILED_TO_SEND_SLACK_MESSAGE = "Failed to send slack message"
)
