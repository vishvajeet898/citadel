package structures

type OmsOrderDetailsResponse struct {
	Id               uint    `json:"id"`
	RequestId        uint    `json:"requestId"`
	Status           uint    `json:"status"`
	SystemDoctorId   uint    `json:"doctorId"`
	PatientName      string  `json:"patientName"`
	PatientNumber    string  `json:"patientNumber"`
	PatientEmail     string  `json:"patientEmail"`
	PatientAge       float32 `json:"patientAge"`
	PatientAgeMonths uint    `json:"patientAgeMonths"`
	PatientAgeDays   uint    `json:"patientAgeDays"`
	PatientGender    string  `json:"patientGender"`
	PatientId        string  `json:"patientId"`
	PatientPK        uint    `json:"patientPk"`
	PartnerId        uint    `json:"partnerId"`
	Notes            string  `json:"notes"`
	ReferredBy       string  `json:"referredBy"`
}
