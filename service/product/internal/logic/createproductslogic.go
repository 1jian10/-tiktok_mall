package logic

import (
	"context"
	"mall/model/database"
	"strconv"

	"mall/service/product/internal/svc"
	"mall/service/product/proto/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateProductsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateProductsLogic {
	return &CreateProductsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateProductsLogic) CreateProducts(in *product.CreateProductsReq) (*product.CreateProductsResp, error) {
	log := l.svcCtx.Log
	db := l.svcCtx.DB
	res := make([]uint32, len(in.Products))
	for i, v := range in.Products {
		p := database.Product{
			Name:        v.Name,
			Description: v.Description,
			Picture:     v.Picture,
			Price:       v.Price,
			Stock:       uint(v.Stock),
		}
		for _, c := range v.Categories {
			p.Categories = append(p.Categories, database.Categories{Name: c})
		}
		if err := db.Model(&database.Product{}).Create(&p).Error; err != nil {
			log.Error(err.Error())
			continue
		}
		log.Info("create product_id:" + strconv.Itoa(int(p.ID)))
		res[i] = uint32(p.ID)
	}
	return &product.CreateProductsResp{ProductId: res}, nil
}
