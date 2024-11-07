package logic

import (
	"context"
	"fmt"
	mlog "mall/log"
	"mall/model/database"

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
	db := l.svcCtx.DB
	u := database.User{}
	err := db.Model(&database.User{}).Where("id = ?", in.UserId).First(&u).Error
	if err != nil {
		mlog.Error(err.Error())
		return &user.InfoResp{Email: ""}, nil
	}
	mlog.Debug("show userinfo:" + fmt.Sprintln(u))
	return &user.InfoResp{Email: u.Email}, nil
}
