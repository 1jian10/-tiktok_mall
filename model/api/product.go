package api

type ListProductsReq struct {
	Page     uint `json:"page"`
	PageSize uint `json:"page_size"`
}
