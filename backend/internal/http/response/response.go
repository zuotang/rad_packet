package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Envelope struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func OK(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, Envelope{Code: "OK", Message: "success", Data: data})
}

func Fail(c echo.Context, status int, code, msg string) error {
	return c.JSON(status, Envelope{Code: code, Message: msg})
}
