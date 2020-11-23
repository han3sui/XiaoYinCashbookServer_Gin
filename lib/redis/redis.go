package redis

import (
	"context"
	"fmt"
	"xiaoyin/lib/config"

	"github.com/go-redis/redis/v8"
)

var Redis *redis.Client
var Nil = redis.Nil

func init() {
	redisConfig := config.Config.GetStringMap("redis")
	Redis = redis.NewClient(&redis.Options{
		Addr:     redisConfig["addr"].(string),
		Password: redisConfig["password"].(string),
		DB:       redisConfig["db"].(int),
	})
	err := Redis.Ping(context.Background()).Err()
	if err != nil {
		panic(fmt.Errorf("redis connect faield: %s \n", err))
	}
}
