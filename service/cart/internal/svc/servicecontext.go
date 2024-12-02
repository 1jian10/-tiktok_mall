package svc

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	mlog "mall/log"
	"mall/service/cart/internal/config"
	"time"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	Log    *mlog.Log
}

func NewServiceContext(c config.Config) *ServiceContext {
	log := mlog.NewLog("CartService")
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

	ctx.DB = db
	ctx.Log = log

	return ctx
}
