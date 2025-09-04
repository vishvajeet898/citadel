package dao

import (
	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/adapters/psql"
	"github.com/Orange-Health/citadel/apps/search/structures"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
)

type SearchDao struct {
	Db *gorm.DB
}

type DataLayer interface {
	// Doctor Tasks
	GetCriticalTaskDetails(taskListDbRequest structures.TaskListDbRequest) (
		[]structures.TaskDetailsDbStruct, *commonStructures.CommonError)
	GetWithheldTaskDetails(taskListDbRequest structures.TaskListDbRequest) (
		[]structures.TaskDetailsDbStruct, *commonStructures.CommonError)
	GetCoAuthorizedTaskDetails(taskListDbRequest structures.TaskListDbRequest) (
		[]structures.TaskDetailsDbStruct, *commonStructures.CommonError)
	GetNormalTaskDetails(taskListDbRequest structures.TaskListDbRequest) (
		[]structures.TaskDetailsDbStruct, *commonStructures.CommonError)
	GetInProgressTaskDetails(taskListDbRequest structures.TaskListDbRequest) (
		[]structures.TaskDetailsDbStruct, *commonStructures.CommonError)
	GetAmendmentTaskDetails(taskListDbRequest structures.TaskListDbRequest) (
		[]structures.TaskDetailsDbStruct, *commonStructures.CommonError)

	// Info Screen
	GetOrderDetailsByBarcode(barcode, serviceType string) (
		structures.InfoScreenOrderDetails, *commonStructures.CommonError)
	GetOrderDetailsByOrderId(omsOrderId string) (
		structures.InfoScreenOrderDetails, *commonStructures.CommonError)
	GetOrderDetailsByVisitId(visitId string) (
		structures.InfoScreenOrderDetails, *commonStructures.CommonError)
	GetVisitBasicDetailsByOrderId(omsOrderId string) (
		[]structures.InfoScreenBasicVisitDetails, *commonStructures.CommonError)
	GetVisitBasicDetailsByBarcode(barcode, serviceType string, labId uint) (
		structures.InfoScreenBasicVisitDetails, *commonStructures.CommonError)
	GetTestDetailsByOrderId(omsOrderId string, servicingCityCode string) (
		[]structures.InfoScreenTestDetails, *commonStructures.CommonError)
	GetTestDetailsByBarcode(barcode, serviceType string, labId uint) (
		[]structures.InfoScreenTestDetails, *commonStructures.CommonError)
}

func InitializeSearchDao() DataLayer {
	return &SearchDao{
		Db: psql.GetDbInstance(),
	}
}
