package util

import (
	"context"
	"github.com/redis/go-redis/v9"
	mlog "mall/log"
	"time"
)

func GetLock(key string, rdb *redis.Client, log *mlog.Log) bool {
	for i := 0; i <= 50; time.Sleep(time.Millisecond * 10) {
		i++
		ok, err := rdb.SetNX(context.Background(), key, "lock", time.Millisecond*100).Result()
		if err != nil {
			log.Warn("get lock:" + err.Error())
			continue
		} else if !ok {
			continue
		}
		return true
	}
	return false
}
