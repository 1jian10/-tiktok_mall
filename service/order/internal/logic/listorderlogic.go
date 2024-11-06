package logic

import (
	"context"

	"mall/service/order/.proto/order"
	"mall/service/order/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOrderLogic {
	return &ListOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListOrderLogic) ListOrder(in *order.ListOrderReq) (*order.ListOrderResp, error) {
	// todo: add your logic here and delete this line

	return &order.ListOrderResp{}, nil
}
