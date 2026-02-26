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

type RewardHandler struct {
	svc *service.RewardService
}

func NewRewardHandler(svc *service.RewardService) *RewardHandler {
	return &RewardHandler{svc: svc}
}

func (h *RewardHandler) Summary(c echo.Context) error {
	data, err := h.svc.Summary(userID(c))
	if err != nil {
		return response.Fail(c, http.StatusInternalServerError, "REWARD_SUMMARY_FAILED", err.Error())
	}
	return response.OK(c, data)
}

func (h *RewardHandler) Unlock(c echo.Context) error {
	count, summary, err := h.svc.UnlockPendingRewards(userID(c))
	if err != nil {
		if errors.Is(err, service.ErrRiskCheckFailed) {
			return response.Fail(c, http.StatusBadRequest, "RISK_CHECK_FAILED", err.Error())
		}
		return response.Fail(c, http.StatusInternalServerError, "REWARD_UNLOCK_FAILED", err.Error())
	}
	return response.OK(c, map[string]interface{}{
		"unlocked_count": count,
		"summary":        summary,
	})
}

func (h *RewardHandler) Records(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	size, _ := strconv.Atoi(c.QueryParam("size"))
	status := strings.TrimSpace(c.QueryParam("status"))
	items, err := h.svc.ListByUser(userID(c), status, page, size)
	if err != nil {
		return response.Fail(c, http.StatusInternalServerError, "REWARD_RECORDS_FAILED", err.Error())
	}
	return response.OK(c, map[string]interface{}{"items": items})
}
