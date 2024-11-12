package logic

import (
	"context"
	"fmt"
	"mall/model/database"
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
			log.Error(err.Error())
			continue
		} else if !ok {
			log.Info("paid get lock false")
			continue
		}
		break
	}
	err := db.Model(&database.Order{}).Where("id = ?", in.OrderId).Update("Paid", "True").Error
	if err != nil {
		log.Error(err.Error())
		rdb.Del(context.Background(), "order:lock"+fmt.Sprintln(in.OrderId))
		return &order.MarkOrderPaidResp{}, err
	}

	err = rdb.Del(context.Background(), "order:lock"+fmt.Sprintln(in.OrderId)).Err()
	if err != nil {
		log.Error(err.Error())
	}
	return &order.MarkOrderPaidResp{}, nil
}
