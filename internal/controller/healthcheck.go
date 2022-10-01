package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type response struct {
	Data string `json:"data"`
}

func HealthCheckHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, response{Data: "Go template at your service!🍲"})
}
