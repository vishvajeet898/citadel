package structures

type PatientDetailsResponse struct {
	Id               uint   `json:"id"`
	Version          int    `json:"version"`
	PatientID        string `json:"patientID"`
	FullName         string `json:"fullName"`
	MobileNumber     string `json:"mobileNumber"`
	PatientEmail     string `json:"patientEmail"`
	SystemCustomerID int    `json:"systemCustomerId"`
	AgeYears         uint   `json:"ageYears"`
	AgeMonths        uint   `json:"ageMonths"`
	AgeDays          uint   `json:"ageDays"`
	Gender           uint   `json:"gender"`
	ContactID        uint   `json:"contactId"`
}
