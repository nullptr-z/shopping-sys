package initialize

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func ViperConfig() (err error) {
	viper.AddConfigPath("user-web/") // 配置文件的路径
	viper.SetConfigName("config")    // 配置文件名
	viper.SetConfigType("yaml")      // 配置文件类型
	if err = viper.ReadInConfig(); err != nil {
		fmt.Println("reading config", err)
		return
	}
	// 热加载配置
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config changes..")
	})
	fmt.Println("Viper Config initialized.......")
	return
	// pg_config := viper.Get("postgres")
	// fmt.Println("GetString:", viper.GetString("postgres.host"))
}
