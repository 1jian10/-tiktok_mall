package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	mlog "mall/log"
	"mall/model"
	"mall/service/order/internal/config"
	"mall/service/order/internal/server"
	"mall/service/order/internal/svc"
	"mall/service/order/proto/order"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/order.yaml", "the config file")

func main() {
	flag.Parse()

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
	go OrderHandle(ctx.RDB, ctx.DB)

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

func OrderHandle(rdb *redis.Client, db *gorm.DB) {
	mlog.Info(fmt.Sprintln("OrderHandler start..."))
	for {
		res, err := rdb.ZRangeByScore(context.Background(), "order:time", &redis.ZRangeBy{
			Min: "0",
			Max: fmt.Sprintf("%f", float64(time.Now().Unix())),
		}).Result()
		if err != nil {
			mlog.Error(err.Error())
			continue
		}
		if len(res) != 0 {
			mlog.Info(fmt.Sprintf("%v:%s", res, "out of time"))
		}
		for _, v := range res {
			ok, err := rdb.SetNX(context.Background(), "order:lock"+v, "lock", time.Millisecond*50).Result()
			if err != nil {
				mlog.Error(err.Error())
				continue
			} else if !ok {
				mlog.Info("delete get lock false")
				continue
			}
			rdb.ZRem(context.Background(), "order:time", v)
			mlog.Info("delete" + fmt.Sprint(v))
			err = db.Where("paid = ?", "False").Delete(&model.Order{}, v).Error
			if err != nil {
				mlog.Error(err.Error())
			}
			err = rdb.Del(context.Background(), "order:lock"+v).Err()
			if err != nil {
				mlog.Error(err.Error())
			}
		}

		time.Sleep(time.Millisecond * 50)
	}

}
