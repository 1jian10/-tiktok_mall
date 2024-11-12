package mlog

import (
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"
	"log/slog"
)

const (
	Debug = 1
	Info  = 2
	Warn  = 3
	Error = 4
)

type Log struct {
	name     string
	producer *nsq.Producer
	level    int
}

type LogBody struct {
	Name    string `json:"name"`
	Level   int    `json:"level"`
	Message string `json:"message"`
}

func NewLog(name string, level ...int) *Log {
	config := nsq.NewConfig()
	p, err := nsq.NewProducer("127.0.0.1:4150", config)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if err := p.Ping(); err != nil {
		slog.Error(err.Error())
	}
	log := &Log{
		name:     name,
		producer: p,
		level:    Debug,
	}
	if len(level) != 0 {
		log.level = level[0]
	}
	return log
}

func (l *Log) Debug(msg string) {
	if l.level <= Debug && l != nil {
		l.send(msg, Debug)
	}
}
func (l *Log) Info(msg string) {
	if l.level <= Info && l != nil {
		l.send(msg, Info)
	}
}
func (l *Log) Warn(msg string) {
	if l.level <= Warn && l != nil {
		l.send(msg, Warn)
	}
}
func (l *Log) Error(msg string) {
	if l.level <= Error && l != nil {
		l.send(msg, Error)
	}
}

func (l *Log) send(msg string, level int) {
	m, _ := json.Marshal(LogBody{
		Name:    "[" + l.name + "]",
		Level:   level,
		Message: msg,
	})
	if err := l.producer.Publish("log", m); err != nil {
		slog.Error(err.Error())
	}
}
