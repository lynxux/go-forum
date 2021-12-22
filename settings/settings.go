package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

//使用结构体保存配置信息
var Conf = new(AppConfig)

type AppConfig struct {
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"`
	Version   string `mapstructure:"version"`
	Port      int    `mapstructure:"port"`
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`

	*LogConfig   `mapstructure:"log"`
	*MySqlConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	FileName   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MySqlConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Port         int    `mapstructure:"port"`
	Password     string `mapstructure:"password"`
	DbName       string `mapstructure:"db_name"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Db       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

func Init() (err error) {
	viper.SetConfigFile("./conf/config.yaml")
	//viper.SetConfigType("yaml") //只用于远程获取配置信息时指定配置文件的类型，非远程时不起作用
	//viper.SetConfigName("config") //用于指定配置文件名，不指定文件后缀，在configpath中查找
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println("viper.ReadInConfig() failed,err: ", err)

		return
		panic(fmt.Errorf("Fatal error config file : %s \n", err))
	}

	//将viper读入的信息反序列化到结构体中，之后通过全局结构体变量Conf访问配置
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Println("viper unmarshal failed! err: ", err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了!")
		//当配置文件修改之后，重新将conf反序列化到全局结构体变量中
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Println("viper unmarshal failed! err: ", err)
		}
	})
	return
}
