package constants

var (
	ERROR_LAB_ID_REQUIRED              = "lab_id is required"
	ERROR_INVALID_SEARCH_TYPE          = "invalid search type"
	ERROR_INVALID_ORDER_ID             = "invalid order_id"
	ERROR_TRF_ID_REQUIRED              = "trf_id is required"
	ERROR_DIGITISATION_ISSUE           = "this is a trf order, please digitize the order first"
	ERROR_SCANNED_AT_INCORRECT_LAB     = "this sample belongs to '%s' lab"
	ERROR_NOT_RECEIVED_REASON_REQUIRED = "not_received_reason is required"
	ERROR_BARCODE_MISMATCH             = "barcode does not match"
)

const (
	SearchTypeOrderId = "order_id"
	SearchTypeBarcode = "barcode"
	SearchTypeTrfId   = "trf_id"
)

var SearchTypes = []string{
	SearchTypeOrderId,
	SearchTypeBarcode,
	SearchTypeTrfId,
}

const (
	LabTypeInhouse   = "inhouse"
	LabTypeInterlab  = "interlab"
	LabTypeOutsource = "outsource"
)
