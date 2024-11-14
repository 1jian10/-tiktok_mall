package svc

import (
	"context"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	mlog "mall/log"
	"mall/service/order/internal/config"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	RDB    *redis.Client
	Log    *mlog.Log
	IsSync bool
}

func NewServiceContext(c config.Config) *ServiceContext {

	log := mlog.NewLog("OrderService")
	ctx := &ServiceContext{Config: c}
	dsn := "root:2549124159f@tcp(127.0.0.1:3306)/mall?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})

	if err != nil {
		ctx.DB = nil
		log.Error(err.Error())
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
		log.Error(err.Error())
		return ctx
	}

	ctx.DB = db
	ctx.RDB = rdb
	ctx.Log = log
	ctx.IsSync = false

	return ctx
}
