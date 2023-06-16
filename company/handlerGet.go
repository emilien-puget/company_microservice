package company

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HandlerGet struct{}

func NewHandlerGet() *HandlerGet {
	return &HandlerGet{}
}

func (g HandlerGet) Handle(c echo.Context) error {
	c.String(http.StatusNotImplemented, "not implemented")
	return nil
}
