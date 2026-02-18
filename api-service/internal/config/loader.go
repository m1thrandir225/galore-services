package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

// BindEnvs automatically binds all mapstructure tags to viper
func BindEnv(v *viper.Viper, cfg any) error {
	t := reflect.TypeOf(cfg)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if tag := field.Tag.Get("mapstructure"); tag != "" {
			if err := v.BindEnv(tag); err != nil {
				return err
			}
		}
	}
	return nil
}

// LoadConfig tries to load a config from a given path
func LoadConfig(path string) (Config, error) {
	var config Config
	v := viper.New()

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // Allow underscores in env variables
	v.AutomaticEnv()                                   // Load environment variables

	env := v.GetString("ENVIRONMENT")
	if env == "" || env == "development" {
		v.AddConfigPath(path) // Path to look for the file
		v.SetConfigFile(".env")
		if err := v.ReadInConfig(); err != nil {
			fmt.Println("No .env file found, relying on environment variables")
		} else {
			fmt.Println("Loaded .env file for local development")
		}
	}

	// Bind environment variables to struct
	//
	if err := BindEnv(v, Config{}); err != nil {
		return config, err
	}
	// viper.BindEnv("POSTGRES_USER")
	// viper.BindEnv("POSTGRES_PASSWORD")
	// viper.BindEnv("POSTGRES_DB")
	// viper.BindEnv("ENVIRONMENT")
	// viper.BindEnv("DB_SOURCE")
	// viper.BindEnv("CACHE_SOURCE")
	// viper.BindEnv("CACHE_PASSWORD")
	// viper.BindEnv("WORKER_SOURCE")
	// viper.BindEnv("TESTING_DB_SOURCE")
	// viper.BindEnv("HTTP_SERVER_ADDRESS")
	// viper.BindEnv("EMBEDDING_SERVER_ADDRESS")
	// viper.BindEnv("EMBEDDING_ACCESS_KEY")
	// viper.BindEnv("MIGRATION_ACCESS_KEY")
	// viper.BindEnv("CATEGORIZER_SERVER_ADDRESS")
	// viper.BindEnv("CATEGORIZER_ACCESS_KEY")
	// viper.BindEnv("TOKEN_SYMMETRIC_KEY")
	// viper.BindEnv("ACCESS_TOKEN_DURATION")
	// viper.BindEnv("REFRESH_TOKEN_DURATION")
	// viper.BindEnv("SMTP_HOST")
	// viper.BindEnv("SMTP_PORT")
	// viper.BindEnv("SMTP_USER")
	// viper.BindEnv("SMTP_PASS")
	// viper.BindEnv("TOTP_SECRET")
	// viper.BindEnv("OPENAI_API_KEY")
	// viper.BindEnv("OPENAI_ASSISTANT_ID")
	// viper.BindEnv("OPENAI_THREAD_URL")
	// viper.BindEnv("STABLE_DIFFUSION_URL")
	// viper.BindEnv("STABLE_DIFFUSION_API_KEY")
	// viper.BindEnv("FIREBASE_SERVICE_KEY")
	//
	if err := v.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}
