package logic

import (
	"context"
	"github.com/redis/go-redis/v9"
	"mall/model"
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
	db := l.svcCtx.DB
	rdb := l.svcCtx.RDB
	log := l.svcCtx.Log
	key := make([]string, 0)
	decr := make([]uint, 0)

	for i, pid := range in.ProductId {
		id := strconv.Itoa(int(pid))
		for ; ; time.Sleep(time.Millisecond * 10) {
			ok, err := rdb.SetNX(context.Background(), "product:lock:"+id, "lock", time.Millisecond*100).Result()
			if err != nil {
				log.Warn("place order get lock:" + err.Error())
				continue
			} else if !ok {
				log.Info("get lock failed")
				continue
			}
			break
		}
		stock, err := rdb.Get(context.Background(), "product:stock:"+id).Result()
		if err == nil {
			log.Info("stock get from redis")
			s, _ := strconv.Atoi(stock)
			if s-int(in.Quantity[i]) < 0 {
				log.Info("stock not enough...rollback")
				rollback(key, decr, rdb, "product:stock:"+id)
				return nil, err
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
			return nil, err
		}
		rdb.Set(context.Background(), "product:stock:"+id, strconv.Itoa(int(product.Stock)-int(in.Quantity[i])), time.Minute*30)
		key = append(key, "product:stock:"+id)
		decr = append(decr, uint(in.Quantity[i]))

		_, err = rdb.Del(context.Background(), "product:lock:"+id).Result()
		if err != nil {
			log.Error("place order del lock:" + err.Error())
		}
	}
	return &order.PlaceOrderResp{}, nil
}

func rollback(key []string, stock []uint, rdb *redis.Client, mutex string) {
	for i := range key {
		rdb.IncrBy(context.Background(), key[i], int64(stock[i]))
	}
	rdb.Del(context.Background(), mutex)
}
