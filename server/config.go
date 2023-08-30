package server

import (
	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo/v4"
	conf "redis-go-dispatcher/config"
)

var (
	redisPool *redis.Pool
	config    conf.Config
)

func StartServer(loadedConfig conf.Config) {
	config = loadedConfig

	redisPool = &redis.Pool{
		MaxIdle:   config.Redis.PoolMaxIdle,
		MaxActive: config.Redis.PoolMaxActive,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(config.Redis.URL)
		},
	}

	e := echo.New()

	BuildRouting(e)

	err := e.Start(":" + config.ServerPort)
	e.Logger.Fatal(err)
}
