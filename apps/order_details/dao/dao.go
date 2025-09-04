package dao

import (
	"gorm.io/gorm"

	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtuls "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

type DataLayer interface {
	GetOrderDetailsByOmsOrderId(omsOrderId string) (commonModels.OrderDetails, *commonStructures.CommonError)
	GetOrderDetailsByOmsRequestId(omsRequestId string) ([]commonModels.OrderDetails, *commonStructures.CommonError)
	GetOrderDetailsByOmsOrderIds(omsOrderIds []string) ([]commonModels.OrderDetails, *commonStructures.CommonError)
	GetOrderDetailsByOmsOrderIdAndServicingLabId(omsOrderId string, servicingLabId uint) (commonModels.OrderDetails,
		*commonStructures.CommonError)
	GetOrderDetailsByTrfId(trfId string) (commonModels.OrderDetails, *commonStructures.CommonError)
	GetOrderDetailsByOmsOrderIdWithTx(tx *gorm.DB, omsOrderId string) (commonModels.OrderDetails,
		*commonStructures.CommonError)

	UpdateOrderDetails(orderDetails commonModels.OrderDetails) (commonModels.OrderDetails, *commonStructures.CommonError)
	CreateOrderDetailsWithTx(tx *gorm.DB, orderDetails commonModels.OrderDetails) (
		commonModels.OrderDetails, *commonStructures.CommonError)
	UpdateOrderDetailsWithTx(tx *gorm.DB, orderDetails commonModels.OrderDetails) (
		commonModels.OrderDetails, *commonStructures.CommonError)
}

func (orderDetailsDao *OrderDetailsDao) GetOrderDetailsByOmsOrderId(omsOrderId string) (
	commonModels.OrderDetails, *commonStructures.CommonError) {

	orderDetail := commonModels.OrderDetails{}
	if err := orderDetailsDao.Db.Find(&orderDetail, "oms_order_id = ?", omsOrderId).Error; err != nil {
		return orderDetail, commonUtuls.HandleORMError(err)
	}

	return orderDetail, nil
}

func (orderDetailsDao *OrderDetailsDao) GetOrderDetailsByOmsRequestId(omsRequestId string) (
	[]commonModels.OrderDetails, *commonStructures.CommonError) {
	orderDetails := []commonModels.OrderDetails{}
	if err := orderDetailsDao.Db.Find(&orderDetails, "oms_request_id = ?", omsRequestId).Error; err != nil {
		return orderDetails, commonUtuls.HandleORMError(err)
	}

	return orderDetails, nil
}

func (orderDetailsDao *OrderDetailsDao) GetOrderDetailsByOmsOrderIds(omsOrderIds []string) (
	[]commonModels.OrderDetails, *commonStructures.CommonError) {
	orderDetails := []commonModels.OrderDetails{}
	if err := orderDetailsDao.Db.Find(&orderDetails, "oms_order_id IN (?)", omsOrderIds).Error; err != nil {
		return orderDetails, commonUtuls.HandleORMError(err)
	}
	return orderDetails, nil
}

func (orderDetailsDao *OrderDetailsDao) GetOrderDetailsByOmsOrderIdAndServicingLabId(omsOrderId string, servicingLabId uint) (
	commonModels.OrderDetails, *commonStructures.CommonError) {

	orderDetail := commonModels.OrderDetails{}
	if err := orderDetailsDao.Db.Find(&orderDetail, "oms_order_id = ? and servicing_lab_id = ?", omsOrderId, servicingLabId).
		Error; err != nil {
		return orderDetail, commonUtuls.HandleORMError(err)
	}

	return orderDetail, nil
}

func (orderDetailsDao *OrderDetailsDao) GetOrderDetailsByOmsOrderIdWithTx(tx *gorm.DB, omsOrderId string) (
	commonModels.OrderDetails, *commonStructures.CommonError) {

	orderDetails := commonModels.OrderDetails{}
	if err := tx.Find(&orderDetails, "oms_order_id = ?", omsOrderId).Error; err != nil {
		return orderDetails, commonUtuls.HandleORMError(err)
	}

	return orderDetails, nil
}

func (orderDetailsDao *OrderDetailsDao) GetOrderDetailsByTrfId(trfId string) (
	commonModels.OrderDetails, *commonStructures.CommonError) {

	orderDetails := commonModels.OrderDetails{}
	if err := orderDetailsDao.Db.Find(&orderDetails, "trf_id = ?", trfId).Error; err != nil {
		return orderDetails, commonUtuls.HandleORMError(err)
	}

	return orderDetails, nil
}

func (orderDetailsDao *OrderDetailsDao) UpdateOrderDetails(orderDetails commonModels.OrderDetails) (
	commonModels.OrderDetails, *commonStructures.CommonError) {

	if err := orderDetailsDao.Db.Save(&orderDetails).Error; err != nil {
		return orderDetails, commonUtuls.HandleORMError(err)
	}

	return orderDetails, nil
}

func (orderDetailsDao *OrderDetailsDao) CreateOrderDetailsWithTx(tx *gorm.DB, orderDetails commonModels.OrderDetails) (
	commonModels.OrderDetails, *commonStructures.CommonError) {

	if err := tx.Create(&orderDetails).Error; err != nil {
		return orderDetails, commonUtuls.HandleORMError(err)
	}

	return orderDetails, nil
}

func (orderDetailsDao *OrderDetailsDao) UpdateOrderDetailsWithTx(tx *gorm.DB, orderDetails commonModels.OrderDetails) (
	commonModels.OrderDetails, *commonStructures.CommonError) {

	if err := tx.Save(&orderDetails).Error; err != nil {
		return orderDetails, commonUtuls.HandleORMError(err)
	}

	return orderDetails, nil
}
