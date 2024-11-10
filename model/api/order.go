package api

type CheckOutReq struct {
	ProductID []uint  `json:"product_id"`
	Quantity  []int32 `json:"quantity"`
	CartID    uint    `json:"cart_id"`
	Address   Address `json:"address"`
	ZipCode   int     `json:"zip_code"`
}

type CheckOutResp struct {
	Success bool `json:"success"`
}

type ChargeReq struct {
	OrderId uint `json:"order_id"`
}

type ChargeResp struct {
	Success bool `json:"success"`
}

type Address struct {
	StreetAddress string `json:"street_address"`
	City          string `json:"city"`
	State         string `json:"state"`
	Country       string `json:"country"`
	ZipCode       string `json:"zip_code"`
}
