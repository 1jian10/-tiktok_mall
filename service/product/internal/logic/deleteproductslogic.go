package logic

import (
	"context"
	"fmt"
	"gorm.io/gorm/clause"
	"mall/model"
	"mall/service/product/internal/svc"
	"mall/service/product/proto/product"
	"math/rand"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteProductsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteProductsLogic {
	return &DeleteProductsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteProductsLogic) DeleteProducts(in *product.DeleteProductsReq) (*product.DeleteProductsResp, error) {
	if !l.svcCtx.IsSync {
		return l.ASyncDelete(in)
	}
	return l.SyncDelete(in)
}

func (l *DeleteProductsLogic) ASyncDelete(in *product.DeleteProductsReq) (*product.DeleteProductsResp, error) {
	db := l.svcCtx.DB
	log := l.svcCtx.Log
	rdb := l.svcCtx.RDB
	res := make([]uint32, len(in.ProductId))

	for i, id := range in.ProductId {
		idstr := strconv.Itoa(int(id))
		for ; ; time.Sleep(time.Millisecond * 10) {
			ok, err := rdb.SetNX(context.Background(), "product:lock:"+idstr, "lock", time.Millisecond*100).Result()
			if err != nil {
				log.Warn("delete product get lock:" + err.Error())
				continue
			} else if !ok {
				log.Info("get lock failed")
				continue
			}
			break
		}
		err := rdb.Set(context.Background(), "product:stock:"+idstr, 0, time.Second*(1800+time.Duration(rand.Int()%100)*10)).Err()
		if err != nil {
			log.Error("delete product form redis:" + err.Error())
			continue
		}
		p := model.Product{}
		tx := db.Begin()
		err = tx.Where("id = ?", id).Clauses(clause.Locking{Strength: "UPDATE"}).First(&p).Error
		if err != nil {
			log.Error("delete product from mysql:" + err.Error())
			tx.Rollback()
			continue
		}
		err = tx.Take(&p).UpdateColumn("stock", 0).Error
		if err != nil {
			log.Error("delete product from mysql:" + err.Error())
			continue
		}
		tx.Commit()
		log.Debug("stock:" + fmt.Sprint(p.Stock))
		res[i] = id
		err = rdb.Del(context.Background(), "product:lock:"+idstr, "lock").Err()
		if err != nil {
			log.Warn("del lock failed:" + err.Error())
		}
	}
	return &product.DeleteProductsResp{ProductId: res}, nil
}

func (l *DeleteProductsLogic) SyncDelete(in *product.DeleteProductsReq) (*product.DeleteProductsResp, error) {
	db := l.svcCtx.DB
	log := l.svcCtx.Log
	res := make([]uint32, len(in.ProductId))

	for i, id := range in.ProductId {
		p := model.Product{}
		tx := db.Begin()
		err := tx.Where("id = ?", id).Clauses(clause.Locking{Strength: "UPDATE"}).Take(&p).Error
		if err != nil {
			log.Error("mysql get lock:" + err.Error())
			tx.Rollback()
			continue
		}
		err = tx.Take(&p).UpdateColumn("stock", 0).Error
		if err != nil {
			log.Error("mysql update stock:" + err.Error())
			tx.Rollback()
			continue
		}
		log.Debug("stock:" + fmt.Sprint(p.Stock))
		tx.Commit()
		res[i] = id
	}
	return &product.DeleteProductsResp{ProductId: res}, nil
}
