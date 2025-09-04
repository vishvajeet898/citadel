package structures

type PatientDetailsUpdateRequest struct {
	Id          uint   `json:"id" binding:"required"`
	Name        string `json:"name"`
	Gender      string `json:"gender"`
	Number      string `json:"number"`
	ExpectedDob string `json:"expected_dob"`
	Dob         string `json:"dob"`
	OrderId     uint   `json:"order_id" binding:"required"`
	CityCode    string `json:"city_code" binding:"required"`
}

type OmsPatientDetailsUpdateRequest struct {
	Name            string `json:"name"`
	Gender          string `json:"gender"`
	Number          string `json:"number"`
	AgeYears        uint   `json:"age_years"`
	AgeMonths       uint   `json:"age_months"`
	AgeDays         uint   `json:"age_days"`
	SystemPatientId string `json:"system_patient_id"`
}
