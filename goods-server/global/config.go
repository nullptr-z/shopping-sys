package global

type MySqlConf struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Dbname   string `mapstructure:"dbname"`
}

type LoggerConf struct {
	Level      string `mapstructure:"level"`
	LogFile    string `mapstructure:"log_file"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type ConsulConf struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Configure struct {
	// 服务本身信息
	Host string `mapstructure:"host"`
	Name string `mapstructure:"name"`
	Mode string

	Mysql  MySqlConf
	Log    LoggerConf
	Consul ConsulConf
}
