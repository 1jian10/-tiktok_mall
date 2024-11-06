package svc

import (
	"mall/service/order/internal/config"

	"context"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	mlog "mall/log"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	RDB    *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {

	mlog.SetName("ProductService")
	ctx := &ServiceContext{Config: c}
	dsn := "root:2549124159f@tcp(127.0.0.1:3306)/mall?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})

	if err != nil {
		ctx.DB = nil
		mlog.Error(err.Error())
		return ctx
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:       "127.0.0.1:6379",
		DB:         0,
		MaxRetries: 1,
	})
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		ctx.RDB = nil
		mlog.Error(err.Error())
		return ctx
	}

	ctx.DB = db
	ctx.RDB = rdb
	mlog.SetName("OrderService")

	return ctx
}
