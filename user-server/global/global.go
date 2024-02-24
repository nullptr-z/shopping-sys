package global

import (
	"github.com/hashicorp/consul/api"
	"gorm.io/gorm"
)

var DB *gorm.DB

var Consul *api.Client

var Conf *Configure
