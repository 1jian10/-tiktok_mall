package logic

import (
	"context"
	"mall/model"
	"mall/service/user/internal/svc"
	"mall/service/user/proto/user"
	"strconv"

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
	log := l.svcCtx.Log
	db := l.svcCtx.DB

	tx := db.Begin()

	err := tx.Where("user_id = ?", in.UserId).Delete(&model.Cart{}).Error
	if err != nil {
		log.Error(err.Error())
		tx.Rollback()
		return nil, err
	}

	err = tx.Delete(&model.User{}, in.UserId).Error
	if err != nil {
		log.Error(err.Error())
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	log.Info("delete user:" + strconv.Itoa(int(in.UserId)))
	return &user.DeleteResp{UserId: in.UserId}, nil
}
