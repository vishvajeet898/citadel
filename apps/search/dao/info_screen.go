package dao

import (
	"net/http"

	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/apps/search/constants"
	"github.com/Orange-Health/citadel/apps/search/structures"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonModels "github.com/Orange-Health/citadel/models"
)

var SelectStringsForOrderDetailsInSearch = []string{
	"order_details.oms_order_id as order_id",
	"order_details.order_status as order_status",
	"order_details.servicing_lab_id as servicing_lab_id",
	"order_details.city_code as servicing_city_code",
	"patient_details.name as patient_name",
	"patient_details.expected_dob as patient_expected_dob",
	"patient_details.dob as patient_dob",
	"patient_details.gender as patient_gender",
	"order_details.oms_request_id as request_id",
	"order_details.camp_id as camp_id",
	"order_details.created_at as created_at",
	"order_details.deleted_at as deleted_at",
}

var SelectStringsForTestDetailsInSearch = []string{
	"test_details.oms_test_id as id",
	"test_details.central_oms_test_id as alnum_test_id",
	"test_details.test_name as test_name",
	"test_details.processing_lab_id as processing_lab_id",
	"test_details.status as status",
	"test_details.department as department",
	"test_details.oms_status as oms_status",
	"test_details.lab_eta as lab_tat",
	"samples.barcode as barcode",
	"test_sample_mapping.sample_number as sample_number",
	"test_sample_mapping.is_rejected as is_sample_rejected",
	"samples.status as sample_status",
	"test_details.created_at as created_at",
	"test_details.deleted_at as deleted_at",
}

var SelectStringsForBasicVisitDetailsInSearch = []string{
	"samples.visit_id as visit_id",
	"samples.barcode as barcode",
	"samples.vial_type_id as vial_type_id",
	"samples.status as current_status",
	"samples.sample_number as sample_number",
	"samples.id as sample_id",
	"samples.parent_sample_id as parent_sample_id",
	"samples.lab_id as lab_id",
	"samples.destination_lab_id as destination_lab_id",
	"samples.created_at as created_at",
	"samples.deleted_at as deleted_at",
	"sample_metadata.collected_at as collected_at",
	"sample_metadata.received_at as received_at",
	"sample_metadata.transferred_at as transferred_at",
	"sample_metadata.outsourced_at as outsourced_at",
	"sample_metadata.rejected_at as rejected_at",
	"sample_metadata.not_received_at as not_received_at",
	"sample_metadata.lis_sync_at as lis_sync_at",
	"sample_metadata.barcode_scanned_at as barcode_scanned_at",
}

func (searchDao *SearchDao) GetOrderDetailsByBarcode(barcode, serviceType string) (structures.InfoScreenOrderDetails, *commonStructures.CommonError) {
	orderDetails := structures.InfoScreenOrderDetails{}

	query := searchDao.Db.Unscoped().Table(commonConstants.TableOrderDetails).
		Select(SelectStringsForOrderDetailsInSearch).
		Joins("INNER JOIN patient_details ON order_details.patient_details_id = patient_details.id").
		Joins("INNER JOIN samples ON samples.oms_order_id = order_details.oms_order_id").
		Where("samples.barcode = ?", barcode)

	if serviceType == constants.ServiceTypeScan {
		query = query.Where("order_details.deleted_at is NULL").
			Where("samples.deleted_at is NULL")
	}

	err := query.Scan(&orderDetails).Error
	if err != nil {
		return orderDetails, &commonStructures.CommonError{
			StatusCode: http.StatusInternalServerError,
			Message:    constants.ERROR_WHILE_FETCHING_ORDER_DETAILS,
		}
	}
	if orderDetails.OrderID == "" {
		return orderDetails, &commonStructures.CommonError{
			StatusCode: http.StatusNotFound,
			Message:    commonConstants.ERROR_BARCODE_NOT_FOUND,
		}
	}

	return orderDetails, nil
}

func (searchDao *SearchDao) GetOrderDetailsByOrderId(omsOrderId string) (structures.InfoScreenOrderDetails, *commonStructures.CommonError) {
	orderDetails := structures.InfoScreenOrderDetails{}

	err := searchDao.Db.Model(&commonModels.OrderDetails{}).
		Select(SelectStringsForOrderDetailsInSearch).
		Joins("INNER JOIN patient_details ON order_details.patient_details_id = patient_details.id").
		Where("order_details.oms_order_id = ?", omsOrderId).Scan(&orderDetails).Error
	if err != nil {
		return orderDetails, &commonStructures.CommonError{
			StatusCode: http.StatusInternalServerError,
			Message:    constants.ERROR_WHILE_FETCHING_ORDER_DETAILS,
		}
	}
	if orderDetails.OrderID == "" {
		return orderDetails, &commonStructures.CommonError{
			StatusCode: http.StatusNotFound,
			Message:    commonConstants.ERROR_ORDER_ID_NOT_FOUND,
		}
	}

	return orderDetails, nil
}

func (searchDao *SearchDao) GetOrderDetailsByVisitId(visitId string) (structures.InfoScreenOrderDetails,
	*commonStructures.CommonError) {
	orderDetails := structures.InfoScreenOrderDetails{}

	err := searchDao.Db.Table("order_details").
		Select(SelectStringsForOrderDetailsInSearch).
		Joins("INNER JOIN patient_details ON order_details.patient_details_id = patient_details.id").
		Joins("INNER JOIN samples ON samples.oms_order_id = order_details.oms_order_id").
		Where("samples.visit_id = ?", visitId).
		Scan(&orderDetails).Error
	if err != nil {
		return orderDetails, &commonStructures.CommonError{
			StatusCode: http.StatusInternalServerError,
			Message:    constants.ERROR_WHILE_FETCHING_ORDER_DETAILS,
		}
	}
	if orderDetails.OrderID == "" {
		return orderDetails, &commonStructures.CommonError{
			StatusCode: http.StatusNotFound,
			Message:    commonConstants.ERROR_VISIT_ID_NOT_FOUND,
		}
	}

	return orderDetails, nil
}

func (searchDao *SearchDao) GetVisitBasicDetailsByOrderId(omsOrderId string) ([]structures.InfoScreenBasicVisitDetails,
	*commonStructures.CommonError) {
	basicVisitDetails := []structures.InfoScreenBasicVisitDetails{}

	err := searchDao.Db.Table("samples").
		Select(SelectStringsForBasicVisitDetailsInSearch).
		Joins("INNER JOIN sample_metadata ON samples.id = sample_metadata.sample_id").
		Where("samples.oms_order_id = ?", omsOrderId).
		Where("samples.deleted_at is NULL").
		Order("samples.id").Scan(&basicVisitDetails).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return basicVisitDetails, &commonStructures.CommonError{
				StatusCode: http.StatusNotFound,
				Message:    commonConstants.ERROR_ORDER_ID_NOT_FOUND,
			}
		}
		return basicVisitDetails, &commonStructures.CommonError{
			StatusCode: http.StatusInternalServerError,
			Message:    constants.ERROR_WHILE_FETCHING_VISIT_DETAILS,
		}
	}

	return basicVisitDetails, nil
}

func (searchDao *SearchDao) GetVisitBasicDetailsByBarcode(barcode, serviceType string, labId uint) (
	structures.InfoScreenBasicVisitDetails, *commonStructures.CommonError) {
	basicVisitDetails := structures.InfoScreenBasicVisitDetails{}

	query := searchDao.Db.Unscoped().Table("samples").
		Select(SelectStringsForBasicVisitDetailsInSearch).
		Joins("INNER JOIN sample_metadata ON samples.id = sample_metadata.sample_id").
		Where("samples.barcode = ?", barcode).
		Where("samples.lab_id = ?", labId)

	if serviceType == constants.ServiceTypeScan {
		query = query.Where("samples.deleted_at is NULL")
	}

	err := query.Order("samples.id").Scan(&basicVisitDetails).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return basicVisitDetails, &commonStructures.CommonError{
				StatusCode: http.StatusNotFound,
				Message:    commonConstants.ERROR_BARCODE_NOT_FOUND,
			}
		}
		return basicVisitDetails, &commonStructures.CommonError{
			StatusCode: http.StatusInternalServerError,
			Message:    constants.ERROR_WHILE_FETCHING_VISIT_DETAILS,
		}
	}

	return basicVisitDetails, nil
}

func (searchDao *SearchDao) GetTestDetailsByOrderId(omsOrderId string, servicingCityCode string) (
	[]structures.InfoScreenTestDetails, *commonStructures.CommonError) {
	testDetails := []structures.InfoScreenTestDetails{}

	err := searchDao.Db.Table(commonConstants.TableTestDetails).
		Select(SelectStringsForTestDetailsInSearch).
		Joins("INNER JOIN test_sample_mapping ON test_details.central_oms_test_id = test_sample_mapping.oms_test_id AND test_sample_mapping.oms_order_id = test_details.oms_order_id").
		Joins("INNER JOIN samples ON test_sample_mapping.sample_number = samples.sample_number AND samples.oms_order_id = test_sample_mapping.oms_order_id").
		Where("test_details.oms_order_id = ?", omsOrderId).
		Where("test_details.city_code = ?", servicingCityCode).
		Where("samples.deleted_at is NULL").
		Where("test_sample_mapping.deleted_at is NULL").
		Order("test_details.central_oms_test_id").
		Scan(&testDetails).Error
	if err != nil {
		return testDetails, &commonStructures.CommonError{
			StatusCode: http.StatusInternalServerError,
			Message:    constants.ERROR_WHILE_FETCHING_TEST_DETAILS,
		}
	}

	return testDetails, nil
}

func (searchDao *SearchDao) GetTestDetailsByBarcode(barcode, serviceType string, labId uint) (
	[]structures.InfoScreenTestDetails, *commonStructures.CommonError) {
	testDetails := []structures.InfoScreenTestDetails{}

	query := searchDao.Db.Unscoped().Table(commonConstants.TableTestDetails).
		Select(SelectStringsForTestDetailsInSearch).
		Joins("INNER JOIN test_sample_mapping ON test_details.central_oms_test_id = test_sample_mapping.oms_test_id").
		Joins("INNER JOIN samples ON test_sample_mapping.sample_number = samples.sample_number AND samples.oms_order_id = test_sample_mapping.oms_order_id").
		Where("samples.barcode = ? AND samples.lab_id = ?", barcode, labId)

	if serviceType == constants.ServiceTypeScan {
		query = query.Where("test_details.deleted_at is NULL").
			Where("samples.deleted_at is NULL")
	}

	err := query.Order("test_details.central_oms_test_id").
		Scan(&testDetails).Error
	if err != nil {
		return testDetails, &commonStructures.CommonError{
			StatusCode: http.StatusInternalServerError,
			Message:    constants.ERROR_WHILE_FETCHING_TEST_DETAILS,
		}
	}

	return testDetails, nil
}
