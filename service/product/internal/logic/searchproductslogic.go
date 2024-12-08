package logic

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"mall/model"
	"mall/service/product/internal/svc"
	"mall/service/product/proto/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchProductsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchProductsLogic {
	return &SearchProductsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchProductsLogic) SearchProducts(in *product.SearchProductsReq) (*product.SearchProductsResp, error) {

	if in.Method == "name" {
		return l.SearchWithName(in)
	}
	return nil, nil
}

func (l *SearchProductsLogic) SearchWithName(in *product.SearchProductsReq) (*product.SearchProductsResp, error) {
	log := l.svcCtx.Log
	db := l.svcCtx.DB
	p := make([]model.Product, 0)

	log.Info("search product name:" + in.Query)
	//此处索引失效
	err := db.Preload("Categories").Where("name LIKE ?", "%"+in.Query+"%").Limit(30).Find(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Info("not found name like:" + in.Query)
		return &product.SearchProductsResp{}, nil
	} else if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	res := &product.SearchProductsResp{
		Results: make([]*product.Product, len(p)),
	}
	for i, v := range p {
		res.Results[i] = &product.Product{
			Id:         uint32(v.ID),
			Name:       v.Name,
			FilePath:   v.FilePath,
			Price:      v.Price,
			ImagePath:  v.ImagePath,
			Categories: make([]string, len(v.Categories)),
		}

		for j, c := range v.Categories {
			res.Results[i].Categories[j] = c.Name
		}
	}

	return res, nil

}
