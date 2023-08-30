package server

import (
	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func BuildRouting(e *echo.Echo) {
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
	result, err := conn.Do("GET", key)
	if err != nil {
		return err
	}

	if result == nil {
		return c.NoContent(http.StatusNotFound)
	}

	data, err := redis.String(result, err)
	if err != nil {
		return err
	}

	return c.JSONBlob(http.StatusOK, []byte(data))
}
