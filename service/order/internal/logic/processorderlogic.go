package logic

import (
	"context"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	mlog "mall/log"
	"mall/model/database"
	"mall/service/order/internal/svc"
	"mall/service/order/proto/order"
	"time"

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
	db := l.svcCtx.DB
	rdb := l.svcCtx.RDB
	var cost float32
	for _, val := range in.OrderItems {
		cost += val.Cost
	}
	o := database.Order{
		Currency: in.UserCurrency,
		Paid:     "False",
		Cost:     cost,
		UserID:   uint(in.UserId),
		Address: &database.Address{
			StreetAddress: in.Address.StreetAddress,
			City:          in.Address.City,
			State:         in.Address.State,
			Country:       in.Address.Country,
			ZipCode:       in.Address.ZipCode,
		},
	}
	tx := db.Begin()
	err := tx.Create(&o).Error
	if err != nil {
		tx.Rollback()
		mlog.Error(err.Error())
		return &order.ProcessOrderResp{}, nil
	}
	for _, val := range in.OrderItems {
		err = tx.Create(&database.OrderProducts{
			OrderID:   o.ID,
			ProductID: uint(val.Item.ProductId),
			Quantity:  uint(val.Item.Quantity),
		}).Error
		if err != nil {
			tx.Rollback()
			mlog.Error(err.Error())
		}
		res := tx.Model(&database.Product{}).Where("id = ?", val.Item.ProductId).UpdateColumn("Stock", gorm.Expr("Stock - ?", val.Item.Quantity))
		if res.Error != nil {
			tx.Rollback()
			mlog.Error(res.Error.Error())
			return &order.ProcessOrderResp{}, nil
		}
	}
	tx.Commit()
	err = rdb.ZAdd(context.Background(), "order:time", redis.Z{
		Score:  float64(time.Now().Add(time.Minute * 15).Unix()),
		Member: o.ID,
	}).Err()
	if err != nil {
		mlog.Error(err.Error())
	}

	return &order.ProcessOrderResp{}, nil
}
