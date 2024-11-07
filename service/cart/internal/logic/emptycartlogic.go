package logic

import (
	"context"
	mlog "mall/log"
	"mall/model"

	"mall/service/cart/internal/svc"
	"mall/service/cart/proto/cart"

	"github.com/zeromicro/go-zero/core/logx"
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
	c := model.Cart{UserID: uint(in.UserId)}
	if err := db.Where("user_id = ?", c.UserID).Take(&c).Error; err != nil {
		mlog.Error(err.Error())
		return &cart.EmptyCartResp{}, nil
	}
	ProductID := make([]uint, len(c.Products))
	for i, v := range c.Products {
		ProductID[i] = v.ID
	}
	err := db.Where("cart_id = ?", c.ID).Where("product_id in ?", ProductID).Delete(&model.CartProducts{}).Error
	if err != nil {
		mlog.Error(err.Error())
	}
	return &cart.EmptyCartResp{}, nil

}