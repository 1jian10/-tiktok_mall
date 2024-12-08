package logic

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"mall/model"
	"mall/service/product/internal/svc"
	"mall/service/product/proto/product"
	"mall/util"
	"strconv"

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
	log := l.svcCtx.Log
	db := l.svcCtx.DB
	rdb := l.svcCtx.RDB
	group := &l.svcCtx.Group
	idstr := strconv.Itoa(int(in.Id))

	str, err := rdb.Get(context.Background(), "product:"+idstr).Result()
	if err == nil {
		log.Info("get product form redis id:" + idstr)
		res := product.Product{}
		if err := json.Unmarshal([]byte(str), &res); err != nil {
			log.Warn("json unmarshal:" + err.Error())
		} else {
			return &product.GetProductResp{Product: &res}, nil
		}
	} else if errors.Is(err, redis.Nil) {
		log.Info("product not in redis id:" + idstr)
	} else {
		log.Warn("get product form redis:" + err.Error())
	}

	//防止缓存击穿
	ans, err, _ := group.Do("id:"+idstr, func() (interface{}, error) {
		p := model.Product{}
		err := db.Preload("Categories").Where("id = ?", in.Id).Take(&p).Error
		return p, err
	})
	p := ans.(model.Product)
	if err != nil {
		log.Error("take product:" + err.Error())
		return nil, err
	}

	res := &product.GetProductResp{
		Product: &product.Product{
			Id:         uint32(p.ID),
			Name:       p.Name,
			FilePath:   p.FilePath,
			ImagePath:  p.ImagePath,
			Price:      p.Price,
			Categories: make([]string, len(p.Categories)),
		},
	}
	for i, v := range p.Categories {
		res.Product.Categories[i] = v.Name
	}

	j, err := json.Marshal(res.Product)
	if err != nil {
		log.Error("json marshal:" + err.Error())
		return nil, err
	}
	err = rdb.Set(context.Background(), "product:"+strconv.Itoa(int(in.Id)), string(j), util.RandTime()).Err()
	if err != nil {
		log.Warn("set product in redis:" + err.Error())
	}
	return res, nil
}
