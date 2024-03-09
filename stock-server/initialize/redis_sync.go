package initialize

import (
	"fmt"
	"stock-server/global"
	. "stock-server/global"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

func ConnRedisSync() error {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", Conf.Redis.Host, Conf.Redis.Port),
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)

	global.Rdsync = redsync.New(pool)

	fmt.Println("Redis sync initialized........")
	return nil
}
