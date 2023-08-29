package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo/v4"
)

type Prefix struct {
	URI         string `yaml:"uri"`
	RedisPrefix string `yaml:"redis_prefix"`
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

var (
	redisPool *redis.Pool
	config    Config
)

func loadConfig(path string) (Config, error) {
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

func main() {
	configPath := "config.yaml"
	loadedConfig, err := loadConfig(configPath)
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	startServer(loadedConfig)
}

func startServer(loadedConfig Config) {
	config = loadedConfig

	redisPool = &redis.Pool{
		MaxIdle:   config.Redis.PoolMaxIdle,
		MaxActive: config.Redis.PoolMaxActive,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(config.Redis.URL)
		},
	}

	e := echo.New()

	buildRouting(e)

	e.Logger.Fatal(e.Start(":" + config.ServerPort))
}

func buildRouting(e *echo.Echo) {
	for _, prefix := range config.Prefixes {

		e.GET(prefix.URI, func(c echo.Context) error {
			return handleGetAll(c, prefix.RedisPrefix)
		})

		e.GET(prefix.URI+"/:id", func(c echo.Context) error {
			return handleGetOne(c, prefix.RedisPrefix)
		})

	}
}

func handleGetAll(c echo.Context, redisPrefix string) error {
	conn := redisPool.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)

	keys, err := redis.Strings(conn.Do("KEYS", redisPrefix+"*"))
	if err != nil {
		return err
	}

	result := strings.Builder{}
	result.WriteString("[")
	for i, key := range keys {
		if i > 0 {
			result.WriteString(",")
		}
		data, err := redis.String(conn.Do("GET", key))
		if err != nil {
			return err
		}

		result.WriteString(data)
	}
	result.WriteString("]")
	return c.JSONBlob(http.StatusOK, []byte(result.String()))
}

func handleGetOne(c echo.Context, redisPrefix string) error {
	id := c.Param("id")

	conn := redisPool.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)

	key := redisPrefix + id
	data, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return err
	}

	return c.JSONBlob(http.StatusOK, []byte(data))
}
