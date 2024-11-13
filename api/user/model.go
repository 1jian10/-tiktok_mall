package user

import (
	"mall/api"
	"mall/service/user/proto/user"
)

type RegisterResp struct {
	Status api.Status `json:"status"`
}

type LoginResp struct {
	Status api.Status `json:"status"`
	Token  string     `json:"token"`
}

type LogOutResp struct {
	Status api.Status `json:"status"`
}

type InfoResp struct {
	Status api.Status     `json:"status"`
	Data   *user.InfoResp `json:"data"`
}

type DeleteResp struct {
	Status api.Status `json:"status"`
}

type UpdateResp struct {
	Status api.Status `json:"status"`
}
