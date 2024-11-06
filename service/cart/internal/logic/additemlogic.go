package logic

import (
	"context"
	"errors"
	"gorm.io/gorm"
	mlog "mall/log"
	"mall/model"
	"strconv"

	"mall/service/cart/internal/svc"
	"mall/service/cart/proto/cart"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddItemLogic {
	return &AddItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddItemLogic) AddItem(in *cart.AddItemReq) (*cart.AddItemResp, error) {
	db := l.svcCtx.DB
	p := model.Product{}
	u := model.User{}

	err := db.Preload("Cart").Where("id = ?", in.UserId).Take(&u).Error
	if err != nil {
		mlog.Error(err.Error())
		return &cart.AddItemResp{}, nil
	}
	err = db.Where("id = ?", in.Item.ProductId).Take(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		mlog.Warn("AddItem:not found product id:" + strconv.Itoa(int(in.Item.ProductId)))
		return &cart.AddItemResp{}, nil
	} else if err != nil {
		mlog.Error(err.Error())
		return &cart.AddItemResp{}, nil
	}

	c := model.CartProducts{CartID: u.Cart.ID, ProductID: uint(in.Item.ProductId), Quantity: uint(in.Item.Quantity)}
	err = db.Save(&c).Error
	if err != nil {
		mlog.Error(err.Error())
		return &cart.AddItemResp{}, nil
	}

	return &cart.AddItemResp{}, nil

}
