package structures

import (
	testDeailsStructures "github.com/Orange-Health/citadel/apps/test_detail/structures"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
)

type UpdateTsmForTestRequest struct {
	TestDetails          testDeailsStructures.TestDetail        `json:"test_details"`
	TsmInfo              commonStructures.TestSampleMappingInfo `json:"tsm_info"`
	VialTypeId           uint                                   `json:"vial_type_id"`
	NewSampleNumber      uint                                   `json:"newSampleNumber,omitempty"`
	ExistingSampleNumber uint                                   `json:"existingSampleNumber"`
	LabId                uint                                   `json:"lab_id"`
	SystemUserId         uint                                   `json:"system_user_id"`
}

type SampleNumberUpdateRequest struct {
	Tests             []TestDetails `json:"tests"`
	CityCode          string        `json:"cityCode"`
	FromProcessingLab bool          `json:"fromProcessingLab"`
}

type TestDetails struct {
	TestID        uint            `json:"id"`
	SampleDetails []SampleDetails `json:"sampleDetails"`
}

type SampleDetails struct {
	VialTypeID           uint `json:"vialTypeId,omitempty"`
	NewSampleNumber      uint `json:"newSampleNumber,omitempty"`
	ExistingSampleNumber uint `json:"existingSampleNumber"`
}
