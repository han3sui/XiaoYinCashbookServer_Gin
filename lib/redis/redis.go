package redis

import (
	"context"
	"fmt"
	"xiaoyin/lib/config"

	"github.com/go-redis/redis/v8"
)

var Redis *redis.Client
var Nil = redis.Nil

func Init() {
	redisConfig := config.Config.GetStringMap("redis")
	Redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s", redisConfig["addr"]),
		Password: fmt.Sprintf("%s", redisConfig["password"]),
		DB:       int(redisConfig["db"].(int64)),
	})
	err := Redis.Ping(context.Background()).Err()
	if err != nil {
		panic(fmt.Errorf("redis连接失败: %s", err))
	}
}
