package internal

import (
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"time"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	DbDriver              string        `mapstructure:"DB_DRIVER"`
	DbSource              string        `mapstructure:"DB_SOURCE"`
	ServerAddress         string        `mapstructure:"SERVER_ADDRESS"`
	AccessTokenTTL        time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	AccessTokenPrivateKey string        `mapstructure:"TOKEN_PRIVATE_KEY"`
	MigrationUrl          string        `mapstructure:"MIGRATION_URL"`
}

func LoadConfig(path string) (*Config, error) {

	//handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})
	logger := slog.New(handler)
	v := viper.NewWithOptions(viper.WithLogger(logger), viper.ExperimentalBindStruct())
	v.AddConfigPath(path)
	v.AutomaticEnv()

	var config Config

	if err := v.ReadInConfig(); err != nil {
		return &config, err
	}

	if err := v.Unmarshal(&config); err != nil {
		return &config, err
	}

	return &config, nil
}
