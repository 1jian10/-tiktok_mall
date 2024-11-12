package logic

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"mall/model/database"
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

	for _, item := range in.Items {
		id := strconv.Itoa(int(item.ProductId))
		for ; ; time.Sleep(time.Millisecond * 50) {
			ok, err := rdb.SetNX(context.Background(), "product:lock:"+id, "lock", time.Millisecond*100).Result()
			if err != nil {
				log.Error(err.Error())
				continue
			} else if !ok {
				log.Info("get lock failed")
				continue
			}
			break
		}
		stock, err := rdb.Get(context.Background(), "product:stock:"+id).Result()
		if !errors.Is(err, redis.Nil) {
			log.Info("stock get from redis")
			s, _ := strconv.Atoi(stock)
			if s-int(item.Quantity) < 0 {
				log.Info("stock not enough...rollback")
				rollback(key, decr, rdb)
				return &order.PlaceOrderResp{Success: "No"}, nil
			}
			rdb.DecrBy(context.Background(), "product:stock:"+id, int64(item.Quantity))
			rdb.Expire(context.Background(), "product:stock:"+id, time.Minute*30)
			key = append(key, "product:stock:"+id)
			decr = append(decr, uint(item.Quantity))
			continue
		}
		p := database.Product{}
		err = db.Select("Stock").Take(&p, item.ProductId).Error
		if err != nil {
			log.Error(err.Error())
			log.Error("rollback")
			rollback(key, decr, rdb)
			return &order.PlaceOrderResp{Success: "No"}, nil
		} else if p.Stock-uint(item.Quantity) < 0 {
			rdb.Set(context.Background(), "product:stock:"+id, strconv.Itoa(int(p.Stock)), time.Minute*30)
			log.Info("stock not enough...rollback")
			rollback(key, decr, rdb)
			return &order.PlaceOrderResp{Success: "No"}, nil
		}
		rdb.Set(context.Background(), "product:stock:"+id, strconv.Itoa(int(p.Stock)-int(item.Quantity)), time.Minute*30)
		key = append(key, "product:stock:"+id)
		decr = append(decr, uint(item.Quantity))

		_, err = rdb.Del(context.Background(), "product:lock:"+id).Result()
		if err != nil {
			log.Error(err.Error())
		}

	}
	return &order.PlaceOrderResp{
		Success: "Yes",
	}, nil

}

func rollback(key []string, stock []uint, rdb *redis.Client) {
	for i := range key {
		rdb.IncrBy(context.Background(), key[i], int64(stock[i]))
	}
}
