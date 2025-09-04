package controller

import (
	"github.com/Orange-Health/citadel/apps/test_detail/service"
)

type TestDetail struct {
	TestDetailService service.TestDetailServiceInterface
}

func InitTestDetailController() *TestDetail {
	return &TestDetail{
		TestDetailService: service.InitializeTestDetailService(),
	}
}
