package redis

import (
	"context"
	"fmt"
	"web_app/settings"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client
var ctx = context.Background()

func Init(redisConfig *settings.RedisConfig) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Password,
		DB:       redisConfig.Db,
		PoolSize: redisConfig.PoolSize,
	})

	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()

	_, err = rdb.Ping(ctx).Result()
	return err
}

func Close() {
	_ = rdb.Close()
}
