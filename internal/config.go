package internal

import (
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	DBSource string `mapstructure:"DB_SOURCE"`
}

func LoadConfig(path string) (*Config, error) {

	viper.AddConfigPath(path)
	//viper.SetConfigFile(".env")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	var config Config

	if err := viper.ReadInConfig(); err != nil {
		return &config, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return &config, err
	}

	return &config, nil
}
