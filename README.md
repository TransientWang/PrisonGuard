PrisonGuard

# 帮助管理 redis 无法直接设置过期时间顶级 key

* 使用 zset 存储需要管理的key 
* key 生成策略自主管理
* 过期时间可控
* 使用方便

初始化
```go
	scoreRedis := RedisKeyManager{
		Rds:        client,
		ManagerKey: "key-manger",
		KeyGenerateStrategy: func(pattern string) (result []*redis.Z) {
			result = append(result, &redis.Z{
				Score:  float64(time.Now().Unix()),
				Member: pattern,
			})
			return
		},
		SurvivalMinute: 1,
		Error:          make(chan error),
	}
```

 在程序中添加需要被管理的key
 
```go
    scoreRedis.AddKey(key)
```

开始任务

```go
    go scoreRedis.BeginWork()
```