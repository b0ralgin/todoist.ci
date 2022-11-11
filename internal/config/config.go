package config

import (
	"errors"
	"path/filepath"

	"github.com/rkoesters/xdg/basedir"
	"github.com/spf13/viper"
)

type Config struct {
	Token string
}

func LoadConfig() (*Config, error) {
	configPath := filepath.Join(basedir.ConfigHome, "todoist") // REQUIRED if the config file does not have the extension in the name
	viper.SetConfigName("config")                              // name of config file (without extension)
	viper.SetConfigType("json")
	viper.AddConfigPath(configPath)
	viper.SetEnvPrefix("todoist")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, ErrConfigFileNotFound
		} else {
			return nil, err
		}
	}
	cfg := &Config{
		Token: viper.GetString("token"),
	}
	return cfg, nil
}

var ErrConfigFileNotFound = errors.New("config not found")
