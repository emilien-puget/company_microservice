package company

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HandlerPatch struct{}

func NewHandlerPatch() *HandlerPatch {
	return &HandlerPatch{}
}

func (hp HandlerPatch) Handle(c echo.Context) error {
	c.String(http.StatusNotImplemented, "not implemented")
	return nil
}
