package logic

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"mall/model/database"
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
	log := l.svcCtx.Log
	p := database.Product{}
	u := database.User{}

	err := db.Preload("Cart").Where("id = ?", in.UserId).Take(&u).Error
	if err != nil {
		log.Error(err.Error())
		return &cart.AddItemResp{}, err
	}
	err = db.Where("id = ?", in.Item.ProductId).Take(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Warn("AddItem:not found product id:" + strconv.Itoa(int(in.Item.ProductId)))
		return &cart.AddItemResp{}, err
	} else if err != nil {
		log.Error(err.Error())
		return &cart.AddItemResp{}, err
	}

	c := database.CartProducts{CartID: u.Cart.ID, ProductID: uint(in.Item.ProductId), Quantity: uint(in.Item.Quantity)}
	err = db.Save(&c).Error
	if err != nil {
		log.Error(err.Error())
		return &cart.AddItemResp{}, err
	}

	return &cart.AddItemResp{}, nil

}
