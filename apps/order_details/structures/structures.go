package structures

import "time"

type OrderDetail struct {
	Id               uint       `json:"id"`
	CreatedAt        *time.Time `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at"`
	CreatedBy        uint       `json:"created_by"`
	UpdatedBy        uint       `json:"updated_by"`
	DeletedBy        uint       `json:"deleted_by"`
	OmsOrderId       string     `json:"oms_order_id"`
	OmsRequestId     string     `json:"oms_request_id"`
	CityCode         string     `json:"oms_city_code"`
	PatientDetailsId uint       `json:"patient_details_id"`
	PartnerId        uint       `json:"partner_id"`
	DoctorId         uint       `json:"doctor_id"`
	TrfId            string     `json:"trf_id"`
	TrfStatus        uint       `json:"trf_status"`
	ServicingLabId   uint       `json:"servicing_lab_id"`
	CollectionType   uint       `json:"collection_type"`
	RequestSource    string     `json:"request_source"`
	BulkOrderId      uint       `json:"bulk_order_id"`
	CampId           uint       `json:"camp_id"`
}
