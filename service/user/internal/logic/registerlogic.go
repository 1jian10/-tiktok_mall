package logic

import (
	"context"
	"mall/model"
	"strconv"

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
	log := l.svcCtx.Log

	res := db.Where("email = ?", in.Email).First(&u)
	if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		log.Info("register:find repeat record")
		return nil, errors.New("repeat record")
	}

	res = db.Create(&model.User{
		Email:    in.Email,
		Password: in.Password,
	})
	if res.Error != nil {
		log.Error(res.Error.Error())
		return nil, res.Error
	}

	db.Where("email = ?", in.Email).First(&u)
	if err := db.Create(&model.Cart{UserID: u.ID}).Error; err != nil {
		log.Error(err.Error())
		return nil, res.Error
	}
	log.Info("register userid:" + strconv.Itoa(int(u.ID)))

	return &user.RegisterResp{UserId: uint32(u.ID)}, nil
}
