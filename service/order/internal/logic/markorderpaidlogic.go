package logic

import (
	"context"
	"fmt"
	"mall/model"
	"time"

	"mall/service/order/internal/svc"
	"mall/service/order/proto/order"

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

	for {
		ok, err := rdb.SetNX(context.Background(), "order:lock"+fmt.Sprintln(in.OrderId), "lock", time.Millisecond*50).Result()
		if err != nil {
			log.Warn("mark order paid get lock:" + err.Error())
			continue
		} else if !ok {
			log.Info("paid get lock false")
			continue
		}
		break
	}
	err := db.Model(&model.Order{}).Where("id = ?", in.OrderId).Update("Paid", "True").Error
	if err != nil {
		log.Error("mark order paid:" + err.Error())
		rdb.Del(context.Background(), "order:lock"+fmt.Sprintln(in.OrderId))
		return nil, err
	}

	err = rdb.Del(context.Background(), "order:lock"+fmt.Sprintln(in.OrderId)).Err()
	if err != nil {
		log.Error("mark order paid del lock:" + err.Error())
	}
	return &order.MarkOrderPaidResp{}, nil
}
