package model

type LogBody struct {
	Name    string `json:"name"`
	Level   int    `json:"level"`
	Message string `json:"message"`
}

const (
	Debug = 1
	Info  = 2
	Warn  = 3
	Error = 4
)
