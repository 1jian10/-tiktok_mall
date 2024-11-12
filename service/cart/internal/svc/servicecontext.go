package svc

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	mlog "mall/log"
	"mall/service/cart/internal/config"
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
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})

	if err != nil {
		ctx.DB = nil
		log.Error(err.Error())
		return ctx
	}

	ctx.DB = db
	ctx.Log = log

	return ctx
}
