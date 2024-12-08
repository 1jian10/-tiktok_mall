package logic

import (
	"context"
	"mall/model"

	"mall/service/user/internal/svc"
	"mall/service/user/proto/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateLogic) Update(in *user.UpdateReq) (*user.UpdateResp, error) {
	db := l.svcCtx.DB
	log := l.svcCtx.Log

	res := db.Model(&model.User{}).Where("id = ?", in.UserId).Update("password", in.Password)
	if res.Error != nil {
		log.Error(res.Error.Error())
		return nil, res.Error
	}
	return &user.UpdateResp{}, nil
}
