package logic

import (
	"context"
	"errors"
	"gorm.io/gorm"
	mlog "mall/log"
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
	db := l.svcCtx.DB
	p := make([]model.Product, 0)
	err := db.Where("name LIKE %?%", in.Query).Find(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		mlog.Info("not found name like:" + in.Query)
		return &product.SearchProductsResp{}, nil
	} else if err != nil {
		mlog.Error(err.Error())
		return &product.SearchProductsResp{}, nil
	}

	res := &product.SearchProductsResp{
		Results: make([]*product.Product, len(p)),
	}
	for i, v := range p {
		res.Results[i] = &product.Product{
			Id:          uint32(v.ID),
			Name:        v.Name,
			Description: v.Description,
			Price:       v.Price,
			Picture:     v.Picture,
			Categories:  make([]string, len(v.Categories)),
		}
		for j, c := range v.Categories {
			res.Results[i].Categories[j] = c.Name
		}
	}

	return res, nil

}