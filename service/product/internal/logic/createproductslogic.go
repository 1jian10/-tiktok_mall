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
		tx := db.Begin()
		p := model.Product{
			Name:       v.Name,
			ImagePath:  v.ImagePath,
			FilePath:   v.FilePath,
			Price:      v.Price,
			Stock:      uint(in.Stock[i]),
			Categories: make([]model.Categories, len(v.Categories)),
		}

		cate := make([]model.Categories, len(v.Categories))
		for j, c := range v.Categories {
			cate[j].Name = c
			//若是cate的名字不存在则创建，若存在则查找id，因为name为index，可以显著加快查找速度
			err := tx.Where("name = ?", c).FirstOrCreate(&cate[j]).Error
			if err != nil {
				log.Error("first or create categories:" + err.Error())
				tx.Rollback()
				continue
			}
			p.Categories[j].ID = cate[j].ID
		}

		if err := tx.Model(&model.Product{}).Create(&p).Error; err != nil {
			log.Error("create product:" + err.Error())
			tx.Rollback()
			continue
		}

		tx.Commit()
		log.Info("create product id:" + strconv.Itoa(int(p.ID)))
		res[i] = uint32(p.ID)
	}
	return &product.CreateProductsResp{ProductId: res}, nil
}
