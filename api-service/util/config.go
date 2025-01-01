package util

import (
	"fmt"
	"log"
	"strings"
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
	MigrationServiceAddress   string        `mapstructure:"MIGRATION_SERVICE_ADDRESS"`
	MigrationServiceKey       string        `mapstructure:"MIGRATION_ACCESS_KEY"`
	CategorizerServiceKey     string        `mapstructure:"CATEGORIZER_ACCESS_KEY"`
	TokenSymmetricKey         string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration       time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration      time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	SMTPHost                  string        `mapstructure:"SMTP_HOST"`
	SMTPPort                  int           `mapstructure:"SMTP_PORT"`
	SMTPUser                  string        `mapstructure:"SMTP_USER"`
	SMTPPass                  string        `mapstructure:"SMTP_PASS"`
	TOTPSecret                string        `mapstructure:"TOTP_SECRET"`
	OpenAIApiKey              string        `mapstructure:"OPENAI_API_KEY"`
	OpenAIAssistantID         string        `mapstructure:"OPENAI_ASSISTANT_ID"`
	OpenAIThreadURL           string        `mapstructure:"OPENAI_THREAD_URL"`
	StableDiffusionURL        string        `mapstructure:"STABLE_DIFFUSION_URL"`
	StableDiffusionApiKey     string        `mapstructure:"STABLE_DIFFUSION_API_KEY"`
	FirebaseServiceKey        string        `mapstructure:"FIREBASE_SERVICE_KEY"`
}

func LoadConfig(path string) (config Config, err error) {
	/*
	 * Load Viper config
	 */
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // Allow underscores in env variables
	viper.AutomaticEnv()                                   // Load environment variables

	env := viper.GetString("ENVIRONMENT")
	if env == "" || env == "development" {
		viper.AddConfigPath(path) // Path to look for the file
		viper.SetConfigFile(".env")
		if err = viper.ReadInConfig(); err != nil {
			fmt.Println("No .env file found, relying on environment variables")
		} else {
			fmt.Println("Loaded .env file for local development")
		}
	}
	// Bind environment variables to struct
	viper.BindEnv("POSTGRES_USER")
	viper.BindEnv("POSTGRES_PASSWORD")
	viper.BindEnv("POSTGRES_DB")
	viper.BindEnv("ENVIRONMENT")
	viper.BindEnv("DB_SOURCE")
	viper.BindEnv("CACHE_SOURCE")
	viper.BindEnv("CACHE_PASSWORD")
	viper.BindEnv("WORKER_SOURCE")
	viper.BindEnv("TESTING_DB_SOURCE")
	viper.BindEnv("HTTP_SERVER_ADDRESS")
	viper.BindEnv("EMBEDDING_SERVER_ADDRESS")
	viper.BindEnv("EMBEDDING_ACCESS_KEY")
	viper.BindEnv("MIGRATION_ACCESS_KEY")
	viper.BindEnv("CATEGORIZER_SERVER_ADDRESS")
	viper.BindEnv("CATEGORIZER_ACCESS_KEY")
	viper.BindEnv("TOKEN_SYMMETRIC_KEY")
	viper.BindEnv("ACCESS_TOKEN_DURATION")
	viper.BindEnv("REFRESH_TOKEN_DURATION")
	viper.BindEnv("SMTP_HOST")
	viper.BindEnv("SMTP_PORT")
	viper.BindEnv("SMTP_USER")
	viper.BindEnv("SMTP_PASS")
	viper.BindEnv("TOTP_SECRET")
	viper.BindEnv("OPENAI_API_KEY")
	viper.BindEnv("OPENAI_ASSISTANT_ID")
	viper.BindEnv("OPENAI_THREAD_URL")
	viper.BindEnv("STABLE_DIFFUSION_URL")
	viper.BindEnv("STABLE_DIFFUSION_API_KEY")
	viper.BindEnv("FIREBASE_SERVICE_KEY")

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	return
}
