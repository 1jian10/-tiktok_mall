package logic

import (
	"context"
	"errors"
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
	//异步下单时禁止修改stock
	if !l.svcCtx.IsSync && len(in.Stock) != 0 {
		return nil, errors.New("ASync make order can not change Stock")
	}
	db := l.svcCtx.DB
	rdb := l.svcCtx.RDB
	log := l.svcCtx.Log
	res := make([]uint32, len(in.Products))

	for i, v := range in.Products {
		p := model.Product{}
		tx := db.Begin()
		//这里开启数据库行级锁，保证并发安全
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
		p.Stock += uint(in.Stock[i])
		cate := make([]model.Categories, len(v.Categories))

		for j, name := range v.Categories {
			cate[j].Name = name
			err = tx.Where("name = ?", name).FirstOrCreate(&cate[j]).Error
			if err != nil {
				break
			}
		}
		if err != nil {
			log.Warn("update get category:" + err.Error())
			tx.Rollback()
			continue
		}

		err = tx.Save(&p).Error
		if err != nil {
			log.Warn("save product:" + err.Error())
			tx.Rollback()
			continue
		}

		err = tx.Model(&p).Association("Categories").Replace(cate)
		if err != nil {
			log.Error("save product:" + err.Error())
			tx.Rollback()
			continue
		}

		tx.Commit()
		//采用先操作数据库后操作redis的方法，一般情况下不会出现问题
		err = rdb.Del(context.Background(), "product:"+strconv.FormatUint(uint64(p.ID), 10)).Err()
		if err != nil {
			log.Error("update product from redis:" + err.Error())
		}
		res[i] = uint32(p.ID)
	}
	return &product.UpdateProductsResp{ProductId: res}, nil
}
