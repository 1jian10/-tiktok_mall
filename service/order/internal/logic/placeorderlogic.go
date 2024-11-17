package logic

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"mall/model"
	"mall/util"
	"strconv"
	"time"

	"mall/service/order/internal/svc"
	"mall/service/order/proto/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type PlaceOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPlaceOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PlaceOrderLogic {
	return &PlaceOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PlaceOrderLogic) PlaceOrder(in *order.PlaceOrderReq) (*order.PlaceOrderResp, error) {
	if l.svcCtx.IsSync {
		return nil, errors.New("you can do this,it is sync to make order")
	}
	db := l.svcCtx.DB
	rdb := l.svcCtx.RDB
	log := l.svcCtx.Log
	key := make([]string, 0)
	decr := make([]uint, 0)

	for i, pid := range in.ProductId {
		id := strconv.Itoa(int(pid))
		if !util.GetLock("product:lock:"+id, rdb, log) {
			rollback(key, decr, rdb, "product:stock:"+id)
			return nil, errors.New("time out")
		}
		stock, err := rdb.Get(context.Background(), "product:stock:"+id).Result()
		if err == nil {
			log.Info("stock get from redis")
			s, _ := strconv.Atoi(stock)
			if s-int(in.Quantity[i]) < 0 {
				log.Info("stock not enough...rollback")
				rollback(key, decr, rdb, "product:stock:"+id)
				return nil, errors.New("stock not enough")
			}
			rdb.DecrBy(context.Background(), "product:stock:"+id, int64(in.Quantity[i]))
			rdb.Expire(context.Background(), "product:stock:"+id, time.Minute*30)
			key = append(key, "product:stock:"+id)
			decr = append(decr, uint(in.Quantity[i]))
			continue
		}

		product := model.Product{}
		err = db.Select("Stock").Take(&product, pid).Error
		if err != nil {
			log.Error("take product id:" + id + ":" + err.Error())
			log.Error("rollback")
			rollback(key, decr, rdb, "product:stock:"+id)
			return nil, err
		} else if product.Stock-uint(in.Quantity[i]) < 0 {
			rdb.Set(context.Background(), "product:stock:"+id, strconv.Itoa(int(product.ID)), time.Minute*30)
			log.Info("stock not enough...rollback")
			rollback(key, decr, rdb, "product:stock:"+id)
			return nil, errors.New("stock not enough")
		}
		rdb.Set(context.Background(), "product:stock:"+id, strconv.Itoa(int(product.Stock)-int(in.Quantity[i])), time.Minute*30)
		key = append(key, "product:stock:"+id)
		decr = append(decr, uint(in.Quantity[i]))

		_, err = rdb.Del(context.Background(), "product:lock:"+id).Result()
		if err != nil {
			log.Error("place order del lock:" + err.Error())
		}
	}
	tx := db.Begin()
	o := model.Order{UserID: uint(in.UserId)}
	if err := tx.Create(&o).Error; err != nil {
		tx.Rollback()
		log.Error("place create order:" + err.Error())
		rollback(key, decr, rdb, "")
		return nil, err
	}
	tx.Commit()
	return &order.PlaceOrderResp{OrderId: uint32(o.ID)}, nil
}

func rollback(key []string, stock []uint, rdb *redis.Client, mutex string) {
	rdb.Del(context.Background(), mutex)
	for i := range key {
		rdb.IncrBy(context.Background(), key[i], int64(stock[i]))
	}
}
