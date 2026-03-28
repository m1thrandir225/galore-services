// Package config handles all the required configuration for the app
package config

import (
	"time"
)

// Config holds all required configuration fields passed by environment variables
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
