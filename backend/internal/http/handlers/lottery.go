package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"red_packet/backend/internal/http/response"
	"red_packet/backend/internal/service"
)

type LotteryHandler struct {
	svc *service.LotteryService
}

func NewLotteryHandler(svc *service.LotteryService) *LotteryHandler {
	return &LotteryHandler{svc: svc}
}

func (h *LotteryHandler) Status(c echo.Context) error {
	status, err := h.svc.GetStatus(userID(c))
	if err != nil {
		return response.Fail(c, http.StatusInternalServerError, "LOTTERY_STATUS_FAILED", err.Error())
	}
	return response.OK(c, status)
}

func (h *LotteryHandler) Spin(c echo.Context) error {
	result, err := h.svc.Spin(userID(c))
	if err != nil {
		if errors.Is(err, service.ErrNoSpinChance) {
			return response.Fail(c, http.StatusBadRequest, "NO_SPIN_CHANCE", err.Error())
		}
		return response.Fail(c, http.StatusInternalServerError, "LOTTERY_SPIN_FAILED", err.Error())
	}
	return response.OK(c, result)
}

func (h *LotteryHandler) Records(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	size, _ := strconv.Atoi(c.QueryParam("size"))
	records, err := h.svc.ListRecords(userID(c), page, size)
	if err != nil {
		return response.Fail(c, http.StatusInternalServerError, "LOTTERY_RECORDS_FAILED", err.Error())
	}
	return response.OK(c, map[string]interface{}{"items": records})
}
