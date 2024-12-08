package logic

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"mall/model"
	"mall/service/order/internal/svc"
	"mall/service/order/proto/order"
	"mall/util"
	"strconv"

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
	group := &l.svcCtx.Group
	// 记录预减的商品和减去的数目，在出错时需要加回去
	key := make([]string, 0)
	decr := make([]uint, 0)

	for i, pid := range in.ProductId {
		id := strconv.Itoa(int(pid))
		//此处存在并发问题，可能查到之后，数据立马过期
		_, err := rdb.Get(context.Background(), "product:stock:"+id).Result()
		if errors.Is(err, redis.Nil) {
			//防止缓存击穿，同一时间只有一个协程能去获取stock
			//在非单机情况下，需要加分布式锁，此处省略
			_, err, _ := group.Do("product:stock:"+id, func() (interface{}, error) {
				product := model.Product{}
				//防止重复获取数据，覆盖正常的stock
				err = db.Select("Stock").Take(&product, pid).Error
				if err == nil {
					err = rdb.Get(context.Background(), "product:stock:"+id).Err()
					if errors.Is(err, redis.Nil) {
						err = nil
						rdb.Set(context.Background(), "product:stock:"+id, product.Stock, util.RandTime())
					}
				}
				return nil, err
			})

			if err != nil {
				log.Error("take product id:" + id + ":" + err.Error())
				log.Error("rollback")
				rollback(key, decr, rdb)
				return nil, err
			}
		} else if err != nil {
			log.Error("get stock from redis:" + err.Error())
			rollback(key, decr, rdb)
			return nil, err
		}
		//此操作为原子操作，无需加锁，只需要判断返回值是否小于0
		res, err := rdb.DecrBy(context.Background(), "product:stock:"+id, int64(in.Quantity[i])).Result()
		rdb.Expire(context.Background(), "product:stock:"+id, util.RandTime())
		key = append(key, "product:stock:"+id)
		decr = append(decr, uint(in.Quantity[i]))
		if res < 0 {
			rollback(key, decr, rdb)
			log.Info("stock not enough...rollback")
			return nil, errors.New("stock not enough")
		}
	}

	o := model.Order{UserID: uint(in.UserId)}
	if err := db.Create(&o).Error; err != nil {
		log.Error("place create order:" + err.Error())
		rollback(key, decr, rdb)
		return nil, err
	}
	return &order.PlaceOrderResp{OrderId: uint32(o.ID)}, nil
}

// 回滚函数，在库存不足或者订单创建失败时调用，将预减的库存加回去
func rollback(key []string, stock []uint, rdb *redis.Client) {
	for i := range key {
		rdb.IncrBy(context.Background(), key[i], int64(stock[i]))
	}
}
