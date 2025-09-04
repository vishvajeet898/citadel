package structures

type MasterVialType struct {
	Id            uint   `json:"id"` // Equivalent to "Container ID" for Attune
	VialType      string `json:"vialType"`
	VialColor     string `json:"vialColor"`
	SampleId      uint   `json:"sampleId"`
	SampleName    string `json:"sampleName"`
	ContainerId   uint   `json:"containerId"`
	ContainerName string `json:"containerName"`
}

type Lab struct {
	Id                   uint                   `json:"id"`
	LabName              string                 `json:"labName,omitempty"`
	Inhouse              bool                   `json:"inhouse,omitempty"`
	City                 string                 `json:"city,omitempty"`
	InterlabTransferMeta []InterlabTransferMeta `json:"interlabTransferMeta,omitempty"`
}

type InterlabTransferMeta struct {
	DestinationLabId uint `json:"destinationLabId"`
	TransferTime     uint `json:"transferTime"`
}

type CdsTestMaster struct {
	Id               uint                     `json:"id"`
	OrangeTestID     uint                     `json:"orangeTestId,omitempty"`
	Department       string                   `json:"department"`
	CollectionsCount uint                     `json:"collectionsCount,omitempty"`
	TestLabMeta      map[uint]CdsTestCityMeta `json:"testLabMeta,omitempty"`
}

type CdsTestCityMeta struct {
	CityCode             string  `json:"cityCode,omitempty"`
	TestCode             string  `json:"testCode,omitempty"`
	LabTat               float32 `json:"labTat"`
	LabId                uint    `json:"labId"`
	InhouseReportEnabled bool    `json:"inhouseReportEnabled"`
}

type CollectionSequenceRequest struct {
	CityCode      string         `json:"city_code"`
	MasterTestIds []uint         `json:"master_test_ids"`
	OrderDetails  []OrderDetails `json:"order_details" binding:"required"`
	CreateSamples bool           `json:"create_samples"`
}

type OrderDetails struct {
	Indentifier           string `json:"identifier"`
	MasterTestIds         []uint `json:"master_test_ids"`
	MasterPackageIds      []uint `json:"master_package_ids"`
	PackageTestIdsInOrder []uint `json:"package_test_ids_in_order"`
}

type CollectionDetails struct {
	Sequence        int             `json:"sequence"`
	MasterTestIds   []uint          `json:"master_test_ids,omitempty"`
	TestVialMapping map[uint][]uint `json:"test_vial_mapping"`
}

type CollectionSequenceResponse struct {
	Collections []CollectionDetails `json:"collections"`
}
