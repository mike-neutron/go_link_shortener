package initializers

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBHost         string `mapstructure:"DB_HOST"`
	DBUserName     string `mapstructure:"DB_USERNAME"`
	DBUserPassword string `mapstructure:"DB_PASSWORD"`
	DBName         string `mapstructure:"DB_DATABASE"`
	DBPort         string `mapstructure:"DB_PORT"`
}

func LoadConfig(path string) (config Config, err error) {

	viper.SetConfigFile(path)
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
