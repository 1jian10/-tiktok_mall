package logic

import (
	"context"

	"mall/service/order/.proto/order"
	"mall/service/order/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProcessOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProcessOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProcessOrderLogic {
	return &ProcessOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProcessOrderLogic) ProcessOrder(in *order.ProcessOrderReq) (*order.ProcessOrderResp, error) {
	// todo: add your logic here and delete this line

	return &order.ProcessOrderResp{}, nil
}
