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

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	db := l.svcCtx.DB
	u := model.User{}

	res := db.Where("email = ?", in.Email).First(&u)
	if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		mlog.Debug("register:find repeat record")
		return &user.RegisterResp{UserId: 0}, nil
	}
	res = db.Create(&model.User{
		Email:    in.Email,
		Password: in.Password,
	})
	if res.Error != nil {
		mlog.Error(res.Error.Error())
		return &user.RegisterResp{UserId: 0}, nil
	}
	db.Where("email = ?", in.Email).First(&u)

	return &user.RegisterResp{UserId: uint32(u.ID)}, nil
}
