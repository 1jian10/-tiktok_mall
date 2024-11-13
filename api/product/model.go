package Product

import (
	"mall/api"
	"mall/service/product/proto/product"
)

type ListResp struct {
	Status api.Status                `json:"status"`
	Data   *product.ListProductsResp `json:"data"`
}

type GetResp struct {
	Status api.Status              `json:"status"`
	Data   *product.GetProductResp `json:"data"`
}

type SearchResp struct {
	Status api.Status                  `json:"status"`
	Data   *product.SearchProductsResp `json:"data"`
}

type CreateResp struct {
	Status api.Status                  `json:"status"`
	Data   *product.CreateProductsResp `json:"data"`
}
