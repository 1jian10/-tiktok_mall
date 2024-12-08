package logic

import (
	"context"
	"mall/model"
	"mall/service/user/internal/svc"
	"mall/service/user/proto/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type InfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoLogic {
	return &InfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *InfoLogic) Info(in *user.InfoReq) (*user.InfoResp, error) {
	log := l.svcCtx.Log
	db := l.svcCtx.DB
	u := model.User{}

	err := db.Model(&model.User{}).Where("id = ?", in.UserId).First(&u).Error
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	//log.Debug("show userinfo:" + fmt.Sprintln(u))
	return &user.InfoResp{Email: u.Email}, nil
}
