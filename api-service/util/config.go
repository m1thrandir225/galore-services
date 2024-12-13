package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Environment               string        `mapstructure:"ENVIRONMENT"`
	DBSource                  string        `mapstructure:"DB_SOURCE"`
	CacheSource               string        `mapstructure:"CACHE_SOURCE"`
	CachePassword             string        `mapstructure:"CACHE_PASSWORD"`
	WorkerSource              string        `mapstructure:"WORKER_SOURCE"`
	TestingDBSource           string        `mapstructure:"TESTING_DB_SOURCE"`
	HTTPServerAddress         string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	EmbeddingServiceAddress   string        `mapstructure:"EMBEDDING_SERVER_ADDRESS"`
	EmbeddingServiceKey       string        `mapstructure:"EMBEDDING_ACCESS_KEY"`
	CategorizerServiceAddress string        `mapstructure:"CATEGORIZER_SERVER_ADDRESS"`
	MigrationServiceKey       string        `mapstructure:"MIGRATION_ACCESS_KEY"`
	CategorizerServiceKey     string        `mapstructure:"CATEGORIZER_ACCESS_KEY"`
	TokenSymmetricKey         string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration       time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration      time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	SMTPHost                  string        `mapstructure:"SMTP_HOST"`
	SMTPPort                  int           `mapstructure:"SMTP_PORT"`
	SMTPUser                  string        `mapstructure:"SMTP_USER"`
	SMTPPass                  string        `mapstructure:"SMTP_PASS"`
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
