package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"mall/model"
	"mall/service/product/internal/svc"
	"mall/service/product/proto/product"
	"mall/util"
	"strconv"
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

// DeleteProducts 暂时不要使用
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
		err := rdb.Set(context.Background(), "product:stock:"+idstr, 0, util.RandTime()).Err()
		if err != nil {
			log.Error("delete product from redis:" + err.Error())
			continue
		}

		p := model.Product{}
		if err := db.Delete(p, in.ProductId).Error; err != nil {
			log.Error("delete product from mysql:" + err.Error())
			continue
		}
		res[i] = id
	}
	return &product.DeleteProductsResp{ProductId: res}, nil
}

func (l *DeleteProductsLogic) SyncDelete(in *product.DeleteProductsReq) (*product.DeleteProductsResp, error) {
	db := l.svcCtx.DB
	log := l.svcCtx.Log
	res := make([]uint32, len(in.ProductId))

	for i, id := range in.ProductId {
		p := model.Product{}
		if err := db.Delete(p, id).Error; err != nil {
			log.Error("delete product from mysql:" + err.Error())
			continue
		}
		res[i] = id
	}
	return &product.DeleteProductsResp{ProductId: res}, nil
}
