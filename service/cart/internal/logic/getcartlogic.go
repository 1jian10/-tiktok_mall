package logic

import (
	"context"
	mlog "mall/log"
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
	err := db.Preload("Cart").First(&u, in.UserId).Error
	if err != nil {
		mlog.Error(err.Error())
		return &cart.GetCartResp{Items: make([]*cart.CartItem, 0)}, nil
	}
	c := make([]model.CartProducts, 0)
	err = db.Model(&model.CartProducts{}).Where("cart_id = ?", u.Cart.ID).Find(&c).Error
	if err != nil {
		mlog.Error(err.Error())
		return &cart.GetCartResp{Items: make([]*cart.CartItem, 0)}, nil
	}
	res := cart.GetCartResp{Items: make([]*cart.CartItem, len(c))}
	for i, v := range c {
		res.Items[i] = new(cart.CartItem)
		res.Items[i].ProductId = uint32(v.ProductID)
		res.Items[i].Quantity = int32(v.Quantity)
	}
	return &res, nil

}
