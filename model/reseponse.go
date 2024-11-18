package model

type Status struct {
	Code     uint   `json:"code"`
	ErrorMsg string `json:"error_msg"`
}

type Response struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data"`
}

const (
	OK         = 200
	ERROR      = 201
	BADREQUEST = 202
	FORBIDDEN  = 203
)
