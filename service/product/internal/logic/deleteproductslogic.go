package logic

import (
	"context"

	"mall/service/product/internal/svc"
	"mall/service/product/proto/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteProductsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteProductsLogic {
	return &DeleteProductsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteProductsLogic) DeleteProducts(in *product.DeleteProductsReq) (*product.DeleteProductsResp, error) {
	// todo: add your logic here and delete this line

	return &product.DeleteProductsResp{}, nil
}
