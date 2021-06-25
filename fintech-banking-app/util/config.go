package util

import "github.com/spf13/viper"

type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") // 可以使用 JSON，XML等其它格式

	// 自动覆盖环境变量值
	viper.AutomaticEnv()

	// 开始读取配置值
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	// 转换为 Config 结构体的变量
	err = viper.Unmarshal(&config)
	return
}
