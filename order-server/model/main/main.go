package main

import (
	"fmt"
	"log"
	"order-server/global"
	"order-server/initialize"
	"order-server/model"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func main() {
	// 初始化。注意顺序不能变
	initialize.ViperConfig()
	initialize.LoadNaocs()

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		global.Conf.Mysql.User,
		global.Conf.Mysql.Password,
		global.Conf.Mysql.Host,
		global.Conf.Mysql.Port,
		global.Conf.Mysql.Dbname,
	)
	fmt.Println("dsn:", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 不要给数据库表添加`s`后缀
		},
	})
	if err != nil {
		panic(err)
	}

	_ = db.AutoMigrate(
		&model.OrderGoods{},
		&model.ShoppingCart{},
		&model.OrderInfo{},
	)
}
