package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf = new(AppConfig)

type AppConfig struct {
	Name         string `mapstructure:"name"`
	Mode         string `mapstructure:"mode"`
	Version      string `mapstructure:"version"`
	Port         int    `mapstructure:"port"`
	StartTime    string `mapstructure:"start_time"`
	MachineID    int64  `mapstructure:"machine_id"`
	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	Password string `mapstructure:"password"`
	PoolSize int    `mapstructure:"pool_size"`
}

func Init(filename string) (err error) {
	// viper.SetConfigName("config")
	// // viper.SetConfigType("yaml")
	// viper.AddConfigPath(".")
	viper.SetConfigFile(filename)
	err = viper.ReadInConfig() // 读取配置信息
	if err != nil {
		fmt.Printf(" viper.ReadInConfig failed, err:%v", err)
		return
	}

	//反序列化到对象中
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Println("viper.Unmarshal failed")
	}

	//改变后自动更新配置
	viper.WatchConfig() // 配置文件的实时watch
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了")
		//反序列化到对象中
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Println("viper.Unmarshal failed")
		}
	})
	return
}
