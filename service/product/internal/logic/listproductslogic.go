package logic

import (
	"context"
	"fmt"
	"mall/model"
	"mall/service/product/internal/svc"
	"mall/service/product/proto/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListProductsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListProductsLogic {
	return &ListProductsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListProductsLogic) ListProducts(in *product.ListProductsReq) (*product.ListProductsResp, error) {
	log := l.svcCtx.Log
	log.Debug("ListReceive:" + fmt.Sprint(in))
	db := l.svcCtx.DB
	p := make([]model.Product, 0)

	err := db.Preload("Categories").Offset(int(in.Page-1) * int(in.PageSize)).Limit(int(in.PageSize)).Find(&p).Error
	if err != nil {
		log.Error("list products:" + err.Error())
		return nil, err
	}

	//log.Debug("ListSearch:" + fmt.Sprint(p))
	res := &product.ListProductsResp{
		Products: make([]*product.Product, len(p)),
	}

	for i, v := range p {
		res.Products[i] = &product.Product{
			Id:         uint32(v.ID),
			Name:       v.Name,
			FilePath:   v.FilePath,
			Price:      v.Price,
			ImagePath:  v.ImagePath,
			Categories: make([]string, len(v.Categories)),
		}

		for j, c := range v.Categories {
			res.Products[i].Categories[j] = c.Name
		}
	}
	return res, nil
}
