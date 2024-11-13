package cart

import (
	"mall/api"
	"mall/service/cart/proto/cart"
)

type AddResp struct {
	Status api.Status `json:"status"`
}

type GetResp struct {
	Status api.Status        `json:"status"`
	Data   *cart.GetCartResp `json:"data"`
}

type EmptyResp struct {
	Status api.Status `json:"status"`
}
