package util

import (
	"math/rand"
	"time"
)

func RandTime() time.Duration {
	return time.Second * (1800 + time.Duration(rand.Int()%100)*10)
}
