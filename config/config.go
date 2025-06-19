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
	HttpClient HttpClient
	Server     Server
	Cache      Cache
}

type HttpClient struct {
	TimeoutSecs int
}

type Server struct {
	Host string
	Port int
}

type Cache struct {
	Dir string
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
	config.Cache.Dir = GetCachingDir(env)
	return nil
}

func GetConfigPath(env string) string {
	if env == envLocal {
		return "../../config/config.yaml"
	}
	return "./config/config.yaml"
}

func GetCachingDir(env string) string {
	if env == envLocal {
		return "../../.cache"
	}
	return "./.cache"
}
