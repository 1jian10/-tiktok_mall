package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/redis/go-redis/v9"
	"mall/model"
	"mall/service/order/internal/config"
	"mall/service/order/internal/server"
	"mall/service/order/internal/svc"
	"mall/service/order/proto/order"
	"mall/util"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/order.yaml", "the config file")

func main() {

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		order.RegisterOrderServiceServer(grpcServer, server.NewOrderServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	go OrderHandle(ctx)

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

// OrderHandle 用于处理超时订单
func OrderHandle(ctx *svc.ServiceContext) {
	log := ctx.Log
	db := ctx.DB
	rdb := ctx.RDB

	log.Info(fmt.Sprintln("OrderHandler start..."))
	for {
		//选取分数低于当前时间戳的所有订单
		res, err := rdb.ZRangeByScore(context.Background(), "order:time", &redis.ZRangeBy{
			Min: "0",
			Max: fmt.Sprintf("%f", float64(time.Now().Unix())),
		}).Result()
		if err != nil {
			log.Error(err.Error())
			continue
		}

		if len(res) != 0 {
			log.Info(fmt.Sprintf("%v:%s", res, "out of time"))
		}

		for _, v := range res {
			//加锁保证支付订单时的并发安全性
			uid, ok := util.GetLock("order:lock"+v, rdb, log)
			if !ok {
				continue
			}

			rdb.ZRem(context.Background(), "order:time", v)
			log.Info("delete" + fmt.Sprint(v))
			err = db.Where("paid = ?", "False").Delete(&model.Order{}, v).Error
			if err != nil {
				log.Error(err.Error())
			}

			util.UnLock("order:lock"+v, rdb, uid)
		}

		time.Sleep(time.Millisecond * 50)
	}

}
