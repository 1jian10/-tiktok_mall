package logic

import (
	"context"
	"mall/model"
	"mall/service/user/internal/svc"
	"mall/service/user/proto/user"

	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	log := l.svcCtx.Log
	db := l.svcCtx.DB
	u := model.User{}

	res := db.Model(&model.User{}).Where("email = ?", in.Email).Where("password = ?", in.Password).First(&u)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		log.Info("login:record not found")
		return nil, res.Error
	} else if res.Error != nil {
		log.Error(res.Error.Error())
		return nil, res.Error
	}
	log.Info("login email:" + in.Email + " " + "password:" + in.Password)

	return &user.LoginResp{UserId: uint32(u.ID)}, nil
}
