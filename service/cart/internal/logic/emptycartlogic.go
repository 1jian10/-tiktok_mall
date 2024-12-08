package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"mall/model"
	"mall/service/cart/internal/svc"
	"mall/service/cart/proto/cart"
)

type EmptyCartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEmptyCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmptyCartLogic {
	return &EmptyCartLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *EmptyCartLogic) EmptyCart(in *cart.EmptyCartReq) (*cart.EmptyCartResp, error) {
	db := l.svcCtx.DB
	log := l.svcCtx.Log
	c := model.Cart{UserID: uint(in.UserId)}

	if err := db.Preload("Products").Where("user_id = ?", c.UserID).Take(&c).Error; err != nil {
		log.Error("get cart with products:" + err.Error())
		return nil, err
	}

	ProductID := make([]uint, len(c.Products))
	for i, v := range c.Products {
		ProductID[i] = v.ID
	}
	//log.Debug("empty cart_id:" + strconv.Itoa(int(c.ID)))

	err := db.Where("cart_id = ?", c.ID).Where("product_id in ?", ProductID).Delete(&model.CartProducts{}).Error
	if err != nil {
		log.Error("empty cart:" + err.Error())
		return nil, err
	}
	return &cart.EmptyCartResp{}, nil

}
