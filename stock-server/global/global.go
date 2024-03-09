package global

import (
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/hashicorp/consul/api"
	"gorm.io/gorm"
)

var DB *gorm.DB

var Consul *api.Client

var Conf *Configure

var Rds *redis.Client

var Rdsync *redsync.Redsync
