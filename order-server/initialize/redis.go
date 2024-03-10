package initialize

import (
	"fmt"
	"order-server/global"
	. "order-server/global"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func ConnRedis() error {
	rds := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", Conf.Redis.Host, Conf.Redis.Port),
		Password: "",
		DB:       0,
		PoolSize: 30,
	})

	if _, err := rds.Ping(rds.Context()).Result(); err != nil {
		zap.L().Error("Redis ping", zap.Error(err))
		panic(fmt.Sprint("Redis connect failed.", err.Error()))
	}
	global.Rds = rds
	fmt.Println("Redis initialized........")
	return nil
}
