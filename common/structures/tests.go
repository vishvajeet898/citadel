package structures

type TestDetailsForLisEvent struct {
	TestId          string `json:"test_id"`
	TestCode        string `json:"test_code"`
	MasterTestId    uint   `json:"master_test_id"`
	MasterPackageId uint   `json:"master_package_id"`
	TestName        string `json:"test_name"`
	TestType        string `json:"test_type"`
	Barcodes        string `json:"barcodes"`
}

type OmsTestDetailsForLis struct {
	OrderId         string `json:"order_id"`
	TestId          string `json:"test_id"`
	TestCode        string `json:"test_code"`
	MasterTestId    uint   `json:"master_test_id"`
	MasterPackageId uint   `json:"master_package_id"`
}
