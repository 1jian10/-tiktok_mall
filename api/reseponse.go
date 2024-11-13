package api

type Status struct {
	Code     uint   `json:"code"`
	ErrorMsg string `json:"error_msg"`
}

const (
	OK         = 200
	ERROR      = 201
	BADREQUEST = 202
	FORBIDDEN  = 203
)
