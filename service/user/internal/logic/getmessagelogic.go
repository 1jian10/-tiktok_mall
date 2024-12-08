package logic

import (
	"context"
	"mall/model"

	"mall/service/user/internal/svc"
	"mall/service/user/proto/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessageLogic {
	return &GetMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMessageLogic) GetMessage(in *user.GetMessageReq) (*user.GetMessageResp, error) {
	db := l.svcCtx.DB
	log := l.svcCtx.Log

	u := model.User{}
	err := db.Preload("Messages").Find(&u, in.UserId).Error
	if err != nil {
		log.Error("take message:" + err.Error())
		return nil, err
	}

	res := &user.GetMessageResp{
		MessageId: make([]uint32, len(u.Messages)),
		Message:   make([]string, len(u.Messages)),
	}

	for i, msg := range u.Messages {
		res.MessageId[i] = uint32(msg.ID)
		res.Message[i] = msg.Message
	}
	return res, nil
}
