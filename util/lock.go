package util

import (
	"context"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	mlog "mall/log"
	"time"
)

func GetLock(key string, rdb *redis.Client, log *mlog.Log) (string, bool) {
	id := uuid.New().String()
	ctx := context.Background()
	for i := 0; i <= 50; time.Sleep(time.Millisecond * 5) {
		i++
		ok, err := rdb.SetNX(ctx, key, id, time.Millisecond*100).Result()
		if err != nil {
			log.Warn("get lock:" + err.Error())
			continue
		} else if !ok {
			continue
		}
		return id, true
	}
	return "", false
}

func UnLock(key string, rdb *redis.Client, id string) {
	ctx := context.Background()
	if rdb.Get(ctx, key).Name() == id {
		rdb.Del(ctx, key)
	}
	return
}
