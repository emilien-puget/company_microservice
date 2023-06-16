package company

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HandlerPost struct{}

func NewHandlerPost() *HandlerPost {
	return &HandlerPost{}
}

func (hp HandlerPost) Handle(c echo.Context) error {
	c.String(http.StatusNotImplemented, "not implemented")
	return nil
}
