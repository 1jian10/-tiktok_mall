package svc

import (
	"context"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	mlog "mall/log"
	"mall/service/order/internal/config"
	"time"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	RDB    *redis.Client
	Log    *mlog.Log
	IsSync bool
	Group  singleflight.Group
}

func NewServiceContext(c config.Config) *ServiceContext {

	log := mlog.NewLog("OrderService")
	ctx := &ServiceContext{Config: c}
	dsn := "root:2549124159f@tcp(127.0.0.1:3306)/mall?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Error(err.Error())
		return ctx
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Warn(err.Error())
	} else {
		sqlDB.SetMaxOpenConns(16)
		sqlDB.SetMaxIdleConns(8)
		sqlDB.SetConnMaxLifetime(time.Hour)
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
