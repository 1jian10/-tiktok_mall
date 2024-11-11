package logic

import (
	"context"
	"fmt"
	mlog "mall/log"
	"mall/model/database"

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
	mlog.Debug("ListReceive:" + fmt.Sprint(in))
	db := l.svcCtx.DB
	p := make([]database.Product, 0)
	err := db.Preload("Categories").Offset(int(in.Page-1) * int(in.PageSize)).Limit(int(in.PageSize)).Find(&p).Error
	if err != nil {
		mlog.Error(err.Error())
		return &product.ListProductsResp{}, nil
	}
	mlog.Debug("ListSearch:" + fmt.Sprint(p))
	res := &product.ListProductsResp{
		Products: make([]*product.Product, len(p)),
	}
	for i, v := range p {
		res.Products[i] = &product.Product{
			Id:          uint32(v.ID),
			Name:        v.Name,
			Description: v.Description,
			Price:       v.Price,
			Picture:     v.Picture,
			Categories:  make([]string, len(v.Categories)),
		}
		for j, c := range v.Categories {
			res.Products[i].Categories[j] = c.Name
		}
	}
	return res, nil
}
