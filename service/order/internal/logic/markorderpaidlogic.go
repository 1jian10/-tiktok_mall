package logic

import (
	"context"
	"errors"
	"fmt"
	"mall/model"
	"mall/service/order/internal/svc"
	"mall/service/order/proto/order"
	"mall/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkOrderPaidLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMarkOrderPaidLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkOrderPaidLogic {
	return &MarkOrderPaidLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MarkOrderPaidLogic) MarkOrderPaid(in *order.MarkOrderPaidReq) (*order.MarkOrderPaidResp, error) {
	rdb := l.svcCtx.RDB
	db := l.svcCtx.DB
	log := l.svcCtx.Log

	if !util.GetLock("order:lock"+fmt.Sprint(in.OrderId), rdb, log) {
		return nil, errors.New("time out")
	}
	err := db.Model(&model.Order{}).Where("id = ?", in.OrderId).Update("Paid", "True").Error
	if err != nil {
		log.Error("mark order paid:" + err.Error())
		rdb.Del(context.Background(), "order:lock"+fmt.Sprint(in.OrderId))
		return nil, err
	}

	err = rdb.Del(context.Background(), "order:lock"+fmt.Sprint(in.OrderId)).Err()
	if err != nil {
		log.Error("mark order paid del lock:" + err.Error())
	}
	return &order.MarkOrderPaidResp{}, nil
}
