package main

import (
	"context"
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
	mlog "mall/log"
	"mall/service/order/proto/order"
)

type MessageHandler struct{}

var log *mlog.Log
var OrderClient order.OrderServiceClient

func (h *MessageHandler) HandleMessage(message *nsq.Message) error {
	req := order.ProcessOrderReq{}
	err := json.Unmarshal(message.Body, &req)
	if err != nil {
		log.Error("json unmarshal:" + err.Error())
		return nil
	}
	_, _ = OrderClient.ProcessOrder(context.Background(), &req)
	return nil
}
func main() {
	conn := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: []string{"127.0.0.1:4379"},
			Key:   "order.rpc",
		},
	})
	OrderClient = order.NewOrderServiceClient(conn.Conn())
	consumer, err := nsq.NewConsumer("order", "process", nsq.NewConfig())
	if err != nil {
		panic(err.Error())
	}
	consumer.AddHandler(&MessageHandler{})
	if err = consumer.ConnectToNSQD("127.0.0.1:4150"); err != nil {
		panic(err.Error())
	}
	log = mlog.NewLog("mq")
	select {}
}
