package mlog

import (
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"log/slog"
	"mall/model/log"
)

var producer *nsq.Producer
var name string

func init() {
	config := nsq.NewConfig()
	p, err := nsq.NewProducer("127.0.0.1:4150", config)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	if err := p.Ping(); err != nil {
		slog.Error(err.Error())
	}
	producer = p
	name = "unknown"
}

func Debug(msg string) {
	send(msg, log.Debug)
}
func Info(msg string) {
	send(msg, log.Info)
}
func Warn(msg string) {
	send(msg, log.Warn)
}
func Error(msg string) {
	send(msg, log.Error)
}

func send(msg string, level int) {
	m, _ := json.Marshal(log.LogBody{
		Name:    name,
		Level:   level,
		Message: msg,
	})
	if err := producer.Publish("log", m); err != nil {
		slog.Error(err.Error())
	}
}

func SetName(n string) {
	name = "[" + n + "]"
}
