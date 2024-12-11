package logic

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	IsSync := l.svcCtx.IsSync
	var cost float32
	o := model.Order{}
	//log.Debug("userid:" + fmt.Sprint(in.UserId))

	if in.OrderId != 0 {
		o.ID = uint(in.OrderId)
		if err := db.First(&o, in.OrderId).Error; err != nil {
			return nil, err
		}
	}

	o.Currency = in.UserCurrency
	o.Paid = "False"
	o.UserID = uint(in.UserId)
	o.Address = &model.Address{
		StreetAddress: in.Address.StreetAddress,
		City:          in.Address.City,
		State:         in.Address.State,
		Country:       in.Address.Country,
		ZipCode:       in.Address.ZipCode,
	}

	tx := db.Begin()
	for _, val := range in.OrderItems {
		p := model.Product{}
		//此处开启行级锁，保证并发安全
		err := tx.Where("id = ?", val.ProductId).Clauses(clause.Locking{Strength: "UPDATE"}).Take(&p).Error
		if err != nil {
			log.Error(err.Error())
			if !IsSync {
				l.notice(in.OrderId, val.ProductId, "search db failed id:"+fmt.Sprint(val.ProductId))
			}
			tx.Rollback()
			return nil, err
		}

		if p.Stock < uint(val.Quantity) {
			log.Info("stock is not enough rollback...")
			if !IsSync {
				l.notice(in.OrderId, in.UserId, "stock not enough id:"+fmt.Sprint(val.ProductId))
			}
			tx.Rollback()
			return nil, err
		}

		err = tx.Model(&p).Where("id = ?", val.ProductId).UpdateColumn("stock", gorm.Expr("stock - ?", val.Quantity)).Error
		if err != nil {
			log.Error("process order get product:" + err.Error())
			if !IsSync {
				l.notice(in.OrderId, in.UserId, "update db failed id:"+fmt.Sprint(val.ProductId))
			}
			tx.Rollback()
			return nil, err
		}
		cost += p.Price * float32(val.Quantity)
	}
	o.Cost = cost
	err := tx.Create(&o).Error
	if err != nil {
		tx.Rollback()
		if !IsSync {
			l.notice(in.OrderId, in.UserId, "save order failed")
		}
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
			log.Error("create order_products:" + err.Error())
			if !IsSync {
				l.notice(in.OrderId, in.UserId, "create join table failed")
			}
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	//将订单加入redis有序结合中，设置15分钟后过期，用于超时取消订单
	err = rdb.ZAdd(context.Background(), "order:time", redis.Z{
		Score:  float64(time.Now().Add(time.Minute * 15).Unix()),
		Member: o.ID,
	}).Err()
	if err != nil {
		log.Error("set order time:" + err.Error())
	}

	return &order.ProcessOrderResp{OrderId: uint32(o.ID)}, nil
}

// 通知函数，异步下单时调用，用于将错误告知用户
func (l *ProcessOrderLogic) notice(orderId uint32, userId uint32, msg string) {
	db := l.svcCtx.DB
	log := l.svcCtx.Log
	message := model.Message{
		Message: "some thing wrong in your order because:" + msg,
		UserID:  uint(userId),
	}
	err := db.Create(&message).Error
	if err != nil {
		log.Error("write message to user:" + err.Error())
	}
	err = db.Delete(&model.Order{}, orderId).Error
	if err != nil {
		log.Error("delete order:" + err.Error())
	}
}
