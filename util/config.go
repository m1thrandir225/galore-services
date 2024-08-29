package util

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Environment         string        `mapstructure:"ENVIRONMENT"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	HTTPServerAddress   string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	/*
	 * Load Viper config
	 */
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}
