package initialize

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func ViperConfig() {
	viper.AddConfigPath(".")      // 配置文件的路径
	viper.SetConfigName("config") // 配置文件名
	viper.SetConfigType("yaml")   // 配置文件类型
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Failed initialized Init ViperConfig", err.Error())
		panic("ViperConfig")
	}
	// 热重载
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config changes..")
	})
	fmt.Println("Viper Config initialized.......")
}
