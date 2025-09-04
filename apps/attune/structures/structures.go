package structures

import (
	commonStructures "github.com/Orange-Health/citadel/common/structures"
)

type AttunePayloadMeta struct {
	PatientSalutation     string
	PatientName           string
	PatientGender         string
	PatientNumber         string
	PatientId             string
	PatientDob            string
	VisitId               string
	SampleCollectedAt     string
	LabId                 uint
	TotalTestsAmount      int
	SrfId                 string
	AttuneTestDetailsList []commonStructures.AttuneTestDetails
	AttuneTestSampleMap   []commonStructures.AttuneTestSampleMap
	TestBarcodeMap        map[string][]commonStructures.AttuneSampleDetails
}
