package order

import (
	"mall/api"
	"mall/service/order/proto/order"
)

type ListResp struct {
	Status api.Status           `json:"status"`
	Data   *order.ListOrderResp `json:"data"`
}

type CheckOutResp struct {
	Status api.Status `json:"status"`
}

type ChargeResp struct {
	Status api.Status `json:"status"`
}
