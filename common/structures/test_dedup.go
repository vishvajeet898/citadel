package structures

type TestDeduplicationRequest struct {
	TestIds    []uint `json:"test_ids"`
	PackageIds []uint `json:"package_ids"`
}

type DedupTestPackageDetail struct {
	DedupTestDetails    []DedupTestDetail    `json:"dedup_test_details,omitempty"`
	DedupPackageDetails []DedupPackageDetail `json:"dedup_package_details,omitempty"`
}

type DedupPackageDetail struct {
	PackageId       uint `json:"package_id,omitempty"`
	MasterPackageId uint `json:"master_package_id,omitempty"`
}

type DedupTestDetail struct {
	TestId          uint `json:"test_id,omitempty"`
	MasterTestId    uint `json:"master_test_id,omitempty"`
	MasterPackageId uint `json:"master_package_id,omitempty"`
}

type TestDeduplicationResponseData struct {
	Data TestDeduplicationResponse `json:"data"`
}

type TestDeduplicationResponse struct {
	Recommendation Recommendation `json:"recommendation"`
	MetaData       MetaData       `json:"meta_data"`
}

type Recommendation struct {
	AddTestIds                   []uint                     `json:"add_test_ids"`
	RemoveTests                  []RemoveTestDetail         `json:"remove_tests"`
	RemovePackages               []RemovePackageDetail      `json:"remove_packages"`
	RemovePackageTests           []RemovePackageTestDetail  `json:"remove_package_tests"`
	PartialOverlapTests          []PartialOverlapTestDetail `json:"partial_overlap_tests"`
	IsManualInterventionRequired bool                       `json:"is_manual_intervention_required"`
}

type PartialOverlapTestDetail struct {
	TestId                uint                   `json:"test_id"`
	PartialOverlapDetails []PartialOverlapDetail `json:"partial_overlap_details"`
}

type PartialOverlapDetail struct {
	OverlappedTestId           uint   `json:"overlapped_test_id"`
	OverLappedInvestigationIds []uint `json:"overlapped_investigation_ids"`
	PackageId                  uint   `json:"package_id"`
}
type OverlapDetail struct {
	OverlappedTestId           uint   `json:"overlapped_test_id"`
	PackageId                  uint   `json:"package_id"`
	OverlappedInvestigationIds []uint `json:"overlapped_investigation_ids,omitempty"`
}

type RemovePackageDetail struct {
	PackageId                 uint   `json:"package_id"`
	OverlappedPackageIds      []uint `json:"overlapped_package_ids,omitempty"`
	IsRepeated                bool   `json:"is_repeated"`
	IsInvestigationOverlapped bool   `json:"is_investigation_overlapped"`
}

type RemoveTestDetail struct {
	TestId                    uint                    `json:"test_id"`
	CompleteOverlapDetails    []CompleteOverlapDetail `json:"complete_overlap_details"`
	PackageIds                []uint                  `json:"package_ids"`
	IsRepeated                bool                    `json:"is_repeated"`
	IsInvestigationOverlapped bool                    `json:"is_investigation_overlapped"`
}

type RemovePackageTestDetail struct {
	TestId                    uint                    `json:"test_id"`
	CompleteOverlapDetails    []CompleteOverlapDetail `json:"complete_overlap_details"`
	PackageIds                []uint                  `json:"package_ids"`
	IsRepeated                bool                    `json:"is_repeated"`
	IsInvestigationOverlapped bool                    `json:"is_investigation_overlapped"`
}

type CompleteOverlapDetail struct {
	OverlappedTestId uint `json:"overlapped_test_id"`
	PackageId        uint `json:"package_id"`
}
type PartialOverlapTest struct {
	TestId                uint            `json:"test_id"`
	PartialOverlapDetails []OverlapDetail `json:"partial_overlap_details,omitempty"`
}

type MetaData struct {
	Package                      Package `json:"package"`
	Test                         Test    `json:"test"`
	IsManualInterventionRequired bool    `json:"is_manual_intervention_required"`
}

type Package struct {
	Test                         Test    `json:"test"`
	Overlap                      Overlap `json:"overlap"`
	IsManualInterventionRequired bool    `json:"is_manual_intervention_required"`
}

type Test struct {
	Overlap         Overlap `json:"overlap"`
	RepeatedTestIds []uint  `json:"repeated_test_ids"`
}

type Overlap struct {
	CompleteOverlap []uint `json:"complete_overlap"`
	PartialOverlap  []uint `json:"partial_overlap"`
}

type DuplicateTestMap struct {
	ID        uint `json:"id"`         // The ID of the duplicate package/test
	IsPackage bool `json:"is_package"` // True if the duplicate is part of a package
}
