package handlers

import (
	"github.com/labstack/echo/v4"

	"red_packet/backend/internal/http/middleware"
)

func userID(c echo.Context) uint {
	v := c.Get(middleware.CtxUserID)
	if v == nil {
		return 0
	}
	return v.(uint)
}
