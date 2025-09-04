package mapper

import (
	"github.com/Orange-Health/citadel/apps/order_details/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

func MapOrderDetailModelToOrderDetail(orderDetails commonModels.OrderDetails) structures.OrderDetail {
	return structures.OrderDetail{
		Id:               orderDetails.Id,
		OmsOrderId:       orderDetails.OmsOrderId,
		OmsRequestId:     orderDetails.OmsRequestId,
		CityCode:         orderDetails.CityCode,
		PatientDetailsId: orderDetails.PatientDetailsId,
		PartnerId:        orderDetails.PartnerId,
		DoctorId:         orderDetails.DoctorId,
		TrfId:            orderDetails.TrfId,
		ServicingLabId:   orderDetails.ServicingLabId,
		CollectionType:   orderDetails.CollectionType,
		RequestSource:    orderDetails.RequestSource,
		BulkOrderId:      orderDetails.BulkOrderId,
		CampId:           orderDetails.CampId,
		CreatedAt:        orderDetails.CreatedAt,
		UpdatedAt:        orderDetails.UpdatedAt,
		DeletedAt:        commonUtils.GetGoLangTimeFromGormDeletedAt(orderDetails.DeletedAt),
		CreatedBy:        orderDetails.CreatedBy,
		UpdatedBy:        orderDetails.UpdatedBy,
		DeletedBy:        orderDetails.DeletedBy,
	}
}
