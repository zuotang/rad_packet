package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"red_packet/backend/internal/http/response"
	"red_packet/backend/internal/service"
)

type ConfigHandler struct {
	svc *service.ConfigService
}

func NewConfigHandler(svc *service.ConfigService) *ConfigHandler {
	return &ConfigHandler{svc: svc}
}

func (h *ConfigHandler) Bootstrap(c echo.Context) error {
	data, err := h.svc.Bootstrap()
	if err != nil {
		return response.Fail(c, http.StatusInternalServerError, "CONFIG_BOOTSTRAP_FAILED", err.Error())
	}
	return response.OK(c, data)
}
