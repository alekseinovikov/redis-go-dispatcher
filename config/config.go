package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Prefix struct {
	URI                  string        `yaml:"uri"`
	RedisPrefix          string        `yaml:"redis_prefix"`
	CacheEnabled         bool          `yaml:"cache_enabled"`
	CacheRefreshDuration time.Duration `yaml:"cache_refresh_duration"`
	CacheTtl             time.Duration `yaml:"cache_ttl"`
}

type RedisConfig struct {
	URL           string `yaml:"url"`
	PoolMaxIdle   int    `yaml:"pool_max_idle"`
	PoolMaxActive int    `yaml:"pool_max_active"`
}
type Config struct {
	ServerPort string      `yaml:"server_port"`
	Redis      RedisConfig `yaml:"redis"`
	Prefixes   []Prefix    `yaml:"prefixes"`
}

func LoadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	config := Config{}
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
