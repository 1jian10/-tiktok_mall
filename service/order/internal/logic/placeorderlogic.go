package logic

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
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
		//这里即使获取数据后立马过期也没关系
		_, err := rdb.Get(context.Background(), "product:stock:"+id).Result()
		if errors.Is(err, redis.Nil) {
			//防止缓存击穿，同一时间只有一个协程能去获取stock
			_, err, _ := group.Do("product:stock:"+id, func() (interface{}, error) {
				product := model.Product{}

				err = db.Select("Stock").Take(&product, pid).Error
				if err != nil {
					return false, err
				}
				//lua脚本保证原子性
				ok, err := util.SetKey(rdb, "product:stock:"+id, product.Stock, util.RandTime())
				if err != nil {
					log.Error("set product stock:" + err.Error())
				}
				return ok, err
			})
			if err != nil {
				rollback(key, decr, rdb)
				return nil, err
			}

		} else if err != nil {
			log.Error("get stock from redis:" + err.Error())
			rollback(key, decr, rdb)
			return nil, err
		}

		ok, err := util.DecrBy(rdb, "product:stock:"+id, int64(in.Quantity[i]), util.RandTime())
		if err != nil {
			rollback(key, decr, rdb)
			log.Error("decr product stock:" + err.Error())
			return nil, err
		} else if !ok {
			rollback(key, decr, rdb)
			log.Info("stock not enough...rollback")
			return nil, errors.New("stock not enough")
		}

		key = append(key, "product:stock:"+id)
		decr = append(decr, uint(in.Quantity[i]))
	}
	//生成唯一订单id，避免数据库操作，可以显著提升系统qps
	id, err := rdb.Incr(context.Background(), "orderid").Result()
	if err != nil {
		rollback(key, decr, rdb)
		log.Error("incr order id:" + err.Error())
		return nil, err
	}
	o := model.Order{
		Model:  gorm.Model{ID: uint(id)},
		UserID: uint(in.UserId),
	}
	//if err := db.Create(&o).Error; err != nil {
	//	log.Error("place create order:" + err.Error())
	//	rollback(key, decr, rdb)
	//	return nil, err
	//}
	return &order.PlaceOrderResp{OrderId: uint32(o.ID)}, nil
}

// 回滚函数，在库存不足或者订单创建失败时调用，将预减的库存加回去
func rollback(key []string, stock []uint, rdb *redis.Client) {
	for i := range key {
		//这里无需判断key是否存在，是因为我们在减库存时重置了其过期时间
		rdb.IncrBy(context.Background(), key[i], int64(stock[i]))
	}
}
