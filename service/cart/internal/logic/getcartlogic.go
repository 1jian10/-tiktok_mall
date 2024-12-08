package logic

import (
	"context"
	"mall/model"
	"mall/service/cart/internal/svc"
	"mall/service/cart/proto/cart"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCartLogic {
	return &GetCartLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCartLogic) GetCart(in *cart.GetCartReq) (*cart.GetCartResp, error) {
	db := l.svcCtx.DB
	u := model.User{}
	log := l.svcCtx.Log

	err := db.Preload("Cart").First(&u, in.UserId).Error
	if err != nil {
		log.Error("get user with cart:" + err.Error())
		return nil, err
	}

	c := make([]model.CartProducts, 0)
	err = db.Model(&model.CartProducts{}).Where("cart_id = ?", u.Cart.ID).Find(&c).Error
	if err != nil {
		log.Error("get cart:" + err.Error())
		return nil, err
	}

	res := cart.GetCartResp{Items: make([]*cart.CartItem, len(c))}
	for i, v := range c {
		res.Items[i] = new(cart.CartItem)
		res.Items[i].ProductId = uint32(v.ProductID)
		res.Items[i].Quantity = uint32(v.Quantity)
	}
	return &res, nil

}
