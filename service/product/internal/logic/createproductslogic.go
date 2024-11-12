package logic

import (
	"context"
	"mall/model"
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
		p := model.Product{
			Name:      v.Name,
			ImagePath: v.ImagePath,
			FilePath:  v.FilePath,
			Price:     v.Price,
			Stock:     uint(in.Stock[i]),
		}
		for _, c := range v.Categories {
			p.Categories = append(p.Categories, model.Categories{Name: c})
		}
		if err := db.Model(&model.Product{}).Create(&p).Error; err != nil {
			log.Error("create product:" + err.Error())
			continue
		}
		log.Info("create product id:" + strconv.Itoa(int(p.ID)))
		res[i] = uint32(p.ID)
	}
	return &product.CreateProductsResp{ProductId: res}, nil
}
