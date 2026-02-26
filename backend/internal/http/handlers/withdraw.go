package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"red_packet/backend/internal/http/response"
	"red_packet/backend/internal/service"
)

type WithdrawHandler struct {
	svc *service.WithdrawService
}

func NewWithdrawHandler(svc *service.WithdrawService) *WithdrawHandler {
	return &WithdrawHandler{svc: svc}
}

func (h *WithdrawHandler) Apply(c echo.Context) error {
	type req struct {
		Amount float64 `json:"amount"`
	}
	var body req
	if err := c.Bind(&body); err != nil {
		return response.Fail(c, http.StatusBadRequest, "BAD_REQUEST", "invalid request")
	}
	data, err := h.svc.Apply(userID(c), body.Amount)
	if err != nil {
		if errors.Is(err, service.ErrInsufficientFunds) {
			return response.Fail(c, http.StatusBadRequest, "INSUFFICIENT_FUNDS", err.Error())
		}
		if errors.Is(err, service.ErrInvalidAmount) {
			return response.Fail(c, http.StatusBadRequest, "INVALID_AMOUNT", err.Error())
		}
		if errors.Is(err, service.ErrWithdrawBelowMin) {
			return response.Fail(c, http.StatusBadRequest, "WITHDRAW_BELOW_MIN", err.Error())
		}
		if errors.Is(err, service.ErrRiskCheckFailed) {
			return response.Fail(c, http.StatusBadRequest, "RISK_CHECK_FAILED", err.Error())
		}
		return response.Fail(c, http.StatusInternalServerError, "WITHDRAW_APPLY_FAILED", err.Error())
	}
	return response.OK(c, data)
}

func (h *WithdrawHandler) Records(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	size, _ := strconv.Atoi(c.QueryParam("size"))
	status := strings.TrimSpace(c.QueryParam("status"))
	items, err := h.svc.ListByUser(userID(c), status, page, size)
	if err != nil {
		return response.Fail(c, http.StatusInternalServerError, "WITHDRAW_RECORDS_FAILED", err.Error())
	}
	return response.OK(c, map[string]interface{}{"items": items})
}
