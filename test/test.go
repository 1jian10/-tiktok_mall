// 性能测试
package main

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
	"mall/service/order/proto/order"
	"math/rand"
	"time"
)

func main() {
	for i := 0; i < 500; i++ {
		req[i] = order.PlaceOrderReq{
			ProductId: []uint32{rand.Uint32()%80 + 31},
			Quantity:  []int32{5},
			UserId:    2,
		}
	}
	for i := 0; i < 128; i++ {
		go run(i)
	}
	time.Sleep(time.Second * 60)

	sum1 := 0
	sum2 := 0
	for i := 0; i < 128; i++ {
		sum1 += errNums[i]
		sum2 += requestNums[i]
	}
	fmt.Println("total errors:", sum1)
	fmt.Println("total requests:", sum2)

}

var errNums [128]int
var requestNums [128]int
var req [500]order.PlaceOrderReq

func run(i int) {
	OrderConn := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: []string{"127.0.0.1:4379"},
			Key:   "order.rpc",
		},
	})
	OrderClient := order.NewOrderServiceClient(OrderConn.Conn())
	client := OrderClient
	fmt.Println("routine:", i, " start")
	t := time.Now().Add(time.Second * 60).Unix()
	ctx := context.Background()

	for time.Now().Unix() < t {
		_, err := client.PlaceOrder(ctx, &req[rand.Uint32()%500])
		if err != nil {
			errNums[i]++
		} else {
			requestNums[i]++
		}
	}
}
