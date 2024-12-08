package util

import (
	"math/rand"
	"time"
)

// RandTime 获取一个随机的时间
func RandTime() time.Duration {
	return time.Second * (1800 + time.Duration(rand.Int()%100)*10)
}
