package company

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HandlerDelete struct{}

func NewHandlerDelete() *HandlerDelete {
	return &HandlerDelete{}
}

func (hp HandlerDelete) Handle(c echo.Context) error {
	c.String(http.StatusNotImplemented, "not implemented")
	return nil
}
