package logic

import (
	"context"
	mlog "mall/log"
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
	db := l.svcCtx.DB
	u := model.User{}

	res := db.Model(&model.User{}).Where("email = ?", in.Email).Where("password = ?", in.Password).First(&u)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		mlog.Info("login:record not found")
		return &user.LoginResp{UserId: 0}, nil
	} else if res.Error != nil {
		mlog.Error(res.Error.Error())
		return &user.LoginResp{UserId: 0}, nil
	}
	mlog.Info("login email:" + in.Email + " " + "password:" + in.Password)

	return &user.LoginResp{UserId: uint32(u.ID)}, nil
}
