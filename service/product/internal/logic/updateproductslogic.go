package logic

import (
	"context"
	"gorm.io/gorm/clause"
	"mall/model"
	"strconv"

	"mall/service/product/internal/svc"
	"mall/service/product/proto/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProductsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProductsLogic {
	return &UpdateProductsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateProductsLogic) UpdateProducts(in *product.UpdateProductsReq) (*product.UpdateProductsResp, error) {
	db := l.svcCtx.DB
	rdb := l.svcCtx.RDB
	log := l.svcCtx.Log
	res := make([]uint32, len(in.Products))
	for i, v := range in.Products {
		p := model.Product{}
		tx := db.Begin()
		err := tx.Where("id = ?", v.Id).Clauses(clause.Locking{Strength: "UPDATE"}).Take(&p).Error
		if err != nil {
			log.Error("update get lock:" + err.Error())
			tx.Rollback()
			continue
		}
		p.FilePath = v.FilePath
		p.ImagePath = v.ImagePath
		p.Price = v.Price
		p.Name = v.Name
		p.Categories = make([]model.Categories, len(v.Categories))
		for j, cate := range v.Categories {
			p.Categories[j].Name = cate
		}
		err = tx.Save(&p).Error
		if err != nil {
			log.Error("save product:" + err.Error())
			tx.Rollback()
			continue
		}
		tx.Commit()
		err = rdb.Del(context.Background(), "product:"+strconv.FormatUint(uint64(p.ID), 10)).Err()
		if err != nil {
			log.Error("update product from redis:" + err.Error())
		}
		res[i] = uint32(p.ID)
	}
	return &product.UpdateProductsResp{ProductId: res}, nil
}
