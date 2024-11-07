package logic

import (
	"context"

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
	// todo: add your logic here and delete this line

	return &product.UpdateProductsResp{}, nil
}
