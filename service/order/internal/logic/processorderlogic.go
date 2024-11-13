package logic

import (
	"context"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"mall/model"
	"mall/service/order/internal/svc"
	"mall/service/order/proto/order"
	"strconv"
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
	log := l.svcCtx.Log
	var cost float32
	o := model.Order{
		Currency: in.UserCurrency,
		Paid:     "False",
		UserID:   uint(in.UserId),
		Address: &model.Address{
			StreetAddress: in.Address.StreetAddress,
			City:          in.Address.City,
			State:         in.Address.State,
			Country:       in.Address.Country,
			ZipCode:       in.Address.ZipCode,
		},
	}
	tx := db.Begin()
	for _, val := range in.OrderItems {
		res := tx.Model(&model.Product{}).Where("id = ?", val.ProductId).UpdateColumn("Stock", gorm.Expr("Stock - ?", val.Quantity))
		if res.Error != nil {
			tx.Rollback()
			log.Error(res.Error.Error())
			return nil, res.Error
		}
		p := model.Product{}
		res = tx.Where("id = ?", val.ProductId).Select("Price").Take(&p)
		if res.Error != nil {
			tx.Rollback()
			log.Error("process order get product:" + res.Error.Error())
			return nil, res.Error
		}
		cost += p.Price * float32(val.Quantity)
	}
	o.Cost = cost
	err := tx.Create(&o).Error
	if err != nil {
		tx.Rollback()
		log.Error("create order" + err.Error())
		return nil, err
	}
	log.Info("create order id:" + strconv.FormatUint(uint64(o.ID), 10))
	for _, val := range in.OrderItems {
		err := tx.Create(&model.OrderProducts{
			OrderID:   o.ID,
			ProductID: uint(val.ProductId),
			Quantity:  uint(val.Quantity),
		}).Error
		if err != nil {
			tx.Rollback()
			log.Error("create order_products:" + err.Error())
		}
	}
	tx.Commit()
	err = rdb.ZAdd(context.Background(), "order:time", redis.Z{
		Score:  float64(time.Now().Add(time.Minute * 15).Unix()),
		Member: o.ID,
	}).Err()
	if err != nil {
		log.Error("set order time:" + err.Error())
	}

	return &order.ProcessOrderResp{}, nil
}
