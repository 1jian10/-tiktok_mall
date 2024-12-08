package logic

import (
	"context"
	"mall/model"
	"strconv"

	"mall/service/order/internal/svc"
	"mall/service/order/proto/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOrderLogic {
	return &ListOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListOrderLogic) ListOrder(in *order.ListOrderReq) (*order.ListOrderResp, error) {
	db := l.svcCtx.DB
	log := l.svcCtx.Log
	u := &model.User{}

	err := db.Preload("Orders.Products").Preload("Orders.Address").Take(&u, in.UserId).Error
	if err != nil {
		log.Error("take User order:" + err.Error())
		return nil, err
	}

	res := &order.ListOrderResp{
		Orders: make([]*order.Order, len(u.Orders)),
	}
	for i, o := range u.Orders {
		res.Orders[i] = &order.Order{
			OrderItems:   make([]*order.OrderItem, len(o.Products)),
			OrderId:      strconv.Itoa(int(o.ID)),
			UserCurrency: o.Currency,
			Address: &order.Address{
				StreetAddress: o.Address.StreetAddress,
				City:          o.Address.City,
				State:         o.Address.State,
				Country:       o.Address.Country,
				ZipCode:       o.Address.ZipCode,
			},
			Email: u.Email,
		}

		for j, p := range o.Products {
			op := model.OrderProducts{}
			err := db.Where("order_id = ?", o.ID).Where("product_id = ?", p.ID).Take(&op).Error
			if err != nil {
				log.Error("search order:" + err.Error())
				return nil, err
			}

			res.Orders[i].OrderItems[j] = &order.OrderItem{
				ProductId: uint32(op.ProductID),
				Quantity:  int32(op.Quantity),
			}
		}
	}

	return res, nil

}
