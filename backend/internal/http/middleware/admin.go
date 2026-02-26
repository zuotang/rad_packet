package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"red_packet/backend/internal/config"
)

func AdminKey(cfg config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := c.Request().Header.Get("X-Admin-Key")
			if key == "" || key != cfg.Admin.Key {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"code":    "ADMIN_UNAUTHORIZED",
					"message": "invalid admin key",
				})
			}
			return next(c)
		}
	}
}
