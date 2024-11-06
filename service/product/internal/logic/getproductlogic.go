package logic

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"log/slog"
	mlog "mall/log"
	"mall/model"
	"strconv"
	"time"

	"mall/service/product/internal/svc"
	"mall/service/product/proto/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductLogic {
	return &GetProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetProductLogic) GetProduct(in *product.GetProductReq) (*product.GetProductResp, error) {
	db := l.svcCtx.DB
	rdb := l.svcCtx.RDB

	str, err := rdb.Get(context.Background(), "product:"+"all:"+strconv.Itoa(int(in.Id))).Result()
	if !errors.Is(err, redis.Nil) {
		res := product.Product{}
		if err := json.Unmarshal([]byte(str), &res); err != nil {
			mlog.Error(err.Error())
		} else {
			return &product.GetProductResp{Product: &res}, nil
		}
	}
	p := model.Product{}
	err = db.Preload("Categories").Where("id = ?", in.Id).Take(&p).Error
	if err != nil {
		slog.Error(err.Error())
		return &product.GetProductResp{Product: &product.Product{Id: 0}}, nil
	}
	res := &product.GetProductResp{
		Product: &product.Product{
			Id:          uint32(p.ID),
			Name:        p.Name,
			Description: p.Description,
			Picture:     p.Picture,
			Price:       p.Price,
			Stock:       uint32(p.Stock),
			Categories:  make([]string, len(p.Categories)),
		},
	}
	for i, v := range p.Categories {
		res.Product.Categories[i] = v.Name
	}
	j, err := json.Marshal(res.Product)
	if err != nil {
		mlog.Warn(err.Error())
		return res, nil
	}
	err = rdb.Set(context.Background(), "product:"+strconv.Itoa(int(in.Id)), string(j), time.Minute*30).Err()
	if err != nil {
		mlog.Warn(err.Error())
	}
	return res, nil
}
