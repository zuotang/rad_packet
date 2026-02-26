package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"red_packet/backend/internal/http/response"
	"red_packet/backend/internal/service"
)

type WalletHandler struct {
	svc *service.WalletService
}

func NewWalletHandler(svc *service.WalletService) *WalletHandler {
	return &WalletHandler{svc: svc}
}

func (h *WalletHandler) Get(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	size, _ := strconv.Atoi(c.QueryParam("size"))
	data, err := h.svc.Get(userID(c), page, size)
	if err != nil {
		return response.Fail(c, http.StatusInternalServerError, "WALLET_FETCH_FAILED", err.Error())
	}
	return response.OK(c, data)
}
