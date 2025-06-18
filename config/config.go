package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/viper"
)

const (
	configEnvKey = "env"
	envLocal     = "local"
)

type Config struct {
	ClientConfig ClientConfig
	ServerConfig ServerConfig
}

type ClientConfig struct {
}

type ServerConfig struct {
	Host string
	Port int
}

var config Config
var once sync.Once

func Get() (Config, error) {
	var err error
	once.Do(func() {
		err = loadConfig()
	})
	if err != nil {
		return config, fmt.Errorf("failed to load config: %v", err)
	}
	return config, nil
}

func loadConfig() error {
	env := os.Getenv(configEnvKey)
	configPath := GetConfigPath(env)

	v := viper.New()
	v.SetConfigFile(configPath)
	v.AddConfigPath(".")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("error while loading config: %v", err)
	}

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return fmt.Errorf("failed to unmarshal config: %v", err)
	}

	config = c
	return nil
}

func GetConfigPath(env string) string {
	if env == envLocal {
		return "../../config/config-local.yaml"
	}
	return "../../config/config.yaml"
}
