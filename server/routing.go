package server

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"redis-go-dispatcher/service"
)

type GetAllService interface {
	GetAll() (string, error)
}

type GetByIdService interface {
	GetById(id string) (string, error)
}

func BuildRouting(e *echo.Echo) {
	for _, prefix := range config.Prefixes {

		jsonService := service.NewJsonService(prefix.RedisPrefix, redisPool)
		e.GET(prefix.URI, func(c echo.Context) error {
			return handleGetAll(c, jsonService)
		})

		e.GET(prefix.URI+"/:id", func(c echo.Context) error {
			return handleGetOne(c, jsonService)
		})

	}
}

func handleGetAll(c echo.Context, service GetAllService) error {
	all, err := service.GetAll()
	if err != nil {
		return err
	}

	return c.JSONBlob(http.StatusOK, []byte(all))
}

func handleGetOne(c echo.Context, service GetByIdService) error {
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
