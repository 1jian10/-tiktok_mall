package logic

import (
	"context"
	mlog "mall/log"
	"mall/model"

	"mall/service/user/internal/svc"
	"mall/service/user/proto/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteLogic) Delete(in *user.DeleteReq) (*user.DeleteResp, error) {
	db := l.svcCtx.DB

	err := db.Model(&model.User{}).Delete(&model.User{}, in.UserId).Error
	if err != nil {
		mlog.Error(err.Error())
		return &user.DeleteResp{UserId: 0}, err
	}
	return &user.DeleteResp{UserId: in.UserId}, nil
}
