package util

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

var decrBy string
var setKey string

func init() {
	decrBy = `
		local key=KEYS[1]
		local decr=ARGV[1]
		local ttl=ARGV[2]
	
		-- 在减库存之前更新ttl，若该key不存在则返回false
		local ok=redis.call("Expire", key,ttl)
		if ok==0 then
			return '0'
		end
	
		-- 预减库存，若库存不足则将库存回滚
		local value=redis.call("DecrBy", key, decr)
		if value<0 then
			redis.call("IncrBy", key,decr)
			return '0'
		end
	
		return '1'
	`

	setKey = `
		local key=KEYS[1]	
		local value=ARGV[1]
		local ttl=ARGV[2]

		-- 在设置前先判断key是否存在，若存在则不设置
		local t=redis.call("TTL",key)
		if t~=-2 then
			return '1'
		end

		redis.call("Set",key,value)
		redis.call("Expire",key,ttl)
		
		return '1'
		`

}

// DecrBy 对key进行减法，若key不存在或者减后小于0则返回false，无错误则返回true
func DecrBy(rdb *redis.Client, key string, decr int64, ttl time.Duration) (bool, error) {
	ok, err := rdb.Eval(context.Background(), decrBy, []string{key}, decr, int(ttl.Seconds())).Result()
	if err != nil {
		return false, err
	}
	return ok.(string) == "1", err

}

// SetKey 设置key的value同时设置ttl
func SetKey(rdb *redis.Client, key string, value uint, ttl time.Duration) (bool, error) {
	_, err := rdb.Eval(context.Background(), setKey, []string{key}, value, int(ttl.Seconds())).Result()

	if err != nil {
		return false, err
	}
	return true, nil
}
