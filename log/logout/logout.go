package main

import (
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"
	"log/slog"
	"mall/model"
	"time"
)

var level int
var mp map[int]string

type MessageHandler struct {
}

func (h *MessageHandler) HandleMessage(message *nsq.Message) error {
	var msg model.LogBody

	err := json.Unmarshal(message.Body, &msg)
	if err != nil {
		slog.Error(err.Error())
		return nil
	}
	output(msg)

	return nil

}
func output(msg model.LogBody) {
	if msg.Level >= level {
		fmt.Println(msg.Name, time.Now().Format("2006-01-02 15:04:05"), mp[msg.Level]+":", msg.Message)
	}
}
func main() {
	consumer, err := nsq.NewConsumer("log", "output", nsq.NewConfig())
	if err != nil {
		panic(err.Error())
	}
	consumer.AddHandler(&MessageHandler{})
	if err = consumer.ConnectToNSQD("127.0.0.1:4150"); err != nil {
		panic(err.Error())
	}
	level = model.Debug
	mp = map[int]string{
		model.Debug: "DEBUG",
		model.Info:  "INFO",
		model.Warn:  "WARN",
		model.Error: "ERROR",
	}
	slog.Info("log start.....")
	select {}
}
