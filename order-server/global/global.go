package global

import (
	"order-server/proto"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/hashicorp/consul/api"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB

	Consul *api.Client

	Conf *Configure

	Rds *redis.Client

	Rdsync *redsync.Redsync

	GoodsSvc proto.GoodsClient

	StockSvc proto.StockClient
)
