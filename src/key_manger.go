package src

import (
	"fmt"

	"time"

	"github.com/go-redis/redis"
)

type RedisKeyMange interface {
	AddKey(key interface{})
	BeginWork()
}
type RedisKeyManager struct {
	Rds                 *redis.Client
	ManagerKey          string
	KeyGenerateStrategy func(pattern string) (result []*redis.Z)
	SurvivalMinute      time.Duration
	Error               chan error
}

func (rds RedisKeyManager) AddKey(pattern interface{}) {
	var parameters []*redis.Z
	switch pattern.(type) {
	case string:
		parameters = append(parameters, rds.KeyGenerateStrategy(pattern.(string))...)
	case []string:
		for _, key := range pattern.([]string) {
			parameters = append(parameters, rds.KeyGenerateStrategy(key)...)
		}
	default:
		panic("pattern invalidatel")
	}
	res := rds.Rds.ZAdd(rds.ManagerKey, parameters...)
	if res != nil && res.Err() != nil {
		rds.Error <- fmt.Errorf("RedisKeyManager add keys err : %s,keys : %v", res.Err(), parameters)
	}
}

func (rds RedisKeyManager) BeginWork() {
	for range time.Tick(rds.SurvivalMinute * time.Minute) {
		rangeRes := rds.Rds.ZRevRangeByScore(rds.ManagerKey, &redis.ZRangeBy{
			Min: "0",
			Max: fmt.Sprintf("%d", time.Now().Unix()-int64(rds.SurvivalMinute)*60),
		})
		if rangeRes != nil && rangeRes.Err() != nil {
			rds.Error <- fmt.Errorf("RedisKeyManager get keys from cache err :%s ", rangeRes.Err())
			continue
		}
		delRes := rds.Rds.Del(rangeRes.Val()...)
		if delRes != nil && delRes.Err() != nil {
			rds.Error <- fmt.Errorf("RedisKeyManager delete keys err :%s ", delRes.Err())
		} else {
			var deleteKeys []interface{}
			for _, v := range rangeRes.Val() {
				deleteKeys = append(deleteKeys, v)
			}
			rds.Rds.ZRem(rds.ManagerKey, deleteKeys...)
		}
	}
}
