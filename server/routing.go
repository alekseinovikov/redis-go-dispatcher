package server

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"redis-go-dispatcher/service"
	"strings"
)

type RedisService interface {
	GetAll() ([]string, error)
	GetById(id string) (string, error)
}

func BuildRouting(e *echo.Echo) {
	for _, prefix := range config.Prefixes {

		jsonService := service.NewJsonService(prefix.RedisPrefix, redisPool)
		var redisService RedisService
		if prefix.CacheEnabled {
			redisService = service.NewCacheService(jsonService, prefix.CacheRefreshDuration, prefix.CacheTtl)
		} else {
			redisService = jsonService
		}

		e.GET(prefix.URI, func(c echo.Context) error {
			return handleGetAll(c, redisService)
		})

		e.GET(prefix.URI+"/:id", func(c echo.Context) error {
			return handleGetOne(c, redisService)
		})

	}
}

func handleGetAll(c echo.Context, service RedisService) error {
	all, err := service.GetAll()
	if err != nil {
		return err
	}

	result := strings.Builder{}
	result.WriteString("[")
	for i, jsonString := range all {
		if i > 0 {
			result.WriteString(",")
		}
		result.WriteString(jsonString)
	}
	result.WriteString("]")

	return c.JSONBlob(http.StatusOK, []byte(result.String()))
}

func handleGetOne(c echo.Context, service RedisService) error {
	id := c.Param("id")
	result, err := service.GetById(id)
	if err != nil {
		return err
	}

	if result == "" {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSONBlob(http.StatusOK, []byte(result))
}
