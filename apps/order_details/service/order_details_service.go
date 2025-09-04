package service

import (
	"net/http"

	"gorm.io/gorm"

	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonModels "github.com/Orange-Health/citadel/models"
)

type OrderDetailsServiceInterface interface {
	GetOrderDetailsByOmsOrderId(omsOrderId string) (commonModels.OrderDetails, *commonStructures.CommonError)
	GetOrderDetailsByOmsRequestId(omsRequestId string) ([]commonModels.OrderDetails, *commonStructures.CommonError)
	GetOrderDetailsByOmsOrderIds(omsOrderIds []string) ([]commonModels.OrderDetails, *commonStructures.CommonError)
	GetOrderDetailsByTrfId(trfId string) (commonModels.OrderDetails, *commonStructures.CommonError)
	GetOrderDetailsByOmsOrderIdWithTx(tx *gorm.DB, omsOrderId string) (commonModels.OrderDetails,
		*commonStructures.CommonError)

	UpdateOrderDetails(orderDetails commonModels.OrderDetails) (commonModels.OrderDetails, *commonStructures.CommonError)
	CreateOrderDetailsWithTx(tx *gorm.DB, orderDetails commonModels.OrderDetails) (commonModels.OrderDetails,
		*commonStructures.CommonError)
	UpdateOrderDetailsWithTx(tx *gorm.DB, orderDetails commonModels.OrderDetails) (commonModels.OrderDetails,
		*commonStructures.CommonError)
}

func (ods *OrderDetailsService) GetOrderDetailsByOmsOrderIdWithTx(tx *gorm.DB, omsOrderId string) (
	commonModels.OrderDetails, *commonStructures.CommonError) {

	orderDetails, cErr := ods.Dao.GetOrderDetailsByOmsOrderIdWithTx(tx, omsOrderId)
	if cErr != nil {
		return orderDetails, cErr
	}

	if orderDetails.Id == 0 {
		return orderDetails, &commonStructures.CommonError{
			Message:    commonConstants.ERROR_ORDER_ID_NOT_FOUND,
			StatusCode: http.StatusNotFound,
		}
	}

	return orderDetails, nil
}

func (ods *OrderDetailsService) GetOrderDetailsByOmsOrderId(omsOrderId string) (
	commonModels.OrderDetails, *commonStructures.CommonError) {

	return ods.Dao.GetOrderDetailsByOmsOrderId(omsOrderId)
}

func (ods *OrderDetailsService) GetOrderDetailsByOmsRequestId(omsRequestId string) (
	[]commonModels.OrderDetails, *commonStructures.CommonError) {
	return ods.Dao.GetOrderDetailsByOmsRequestId(omsRequestId)
}

func (ods *OrderDetailsService) GetOrderDetailsByOmsOrderIds(omsOrderIds []string) (
	[]commonModels.OrderDetails, *commonStructures.CommonError) {
	return ods.Dao.GetOrderDetailsByOmsOrderIds(omsOrderIds)
}

func (ods *OrderDetailsService) GetOrderDetailsByTrfId(trfId string) (
	commonModels.OrderDetails, *commonStructures.CommonError) {

	return ods.Dao.GetOrderDetailsByTrfId(trfId)
}

func (ods *OrderDetailsService) UpdateOrderDetails(orderDetails commonModels.OrderDetails) (
	commonModels.OrderDetails, *commonStructures.CommonError) {

	return ods.Dao.UpdateOrderDetails(orderDetails)
}

func (ods *OrderDetailsService) CreateOrderDetailsWithTx(tx *gorm.DB,
	orderDetails commonModels.OrderDetails) (commonModels.OrderDetails, *commonStructures.CommonError) {

	return ods.Dao.CreateOrderDetailsWithTx(tx, orderDetails)
}

func (ods *OrderDetailsService) UpdateOrderDetailsWithTx(tx *gorm.DB,
	orderDetails commonModels.OrderDetails) (commonModels.OrderDetails, *commonStructures.CommonError) {

	return ods.Dao.UpdateOrderDetailsWithTx(tx, orderDetails)
}
