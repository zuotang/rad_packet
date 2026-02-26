package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"red_packet/backend/internal/http/response"
	"red_packet/backend/internal/service"
)

type ReferralHandler struct {
	svc *service.ReferralService
}

func NewReferralHandler(svc *service.ReferralService) *ReferralHandler {
	return &ReferralHandler{svc: svc}
}

func (h *ReferralHandler) Bind(c echo.Context) error {
	type req struct {
		Code string `json:"code"`
	}
	var body req
	if err := c.Bind(&body); err != nil || body.Code == "" {
		return response.Fail(c, http.StatusBadRequest, "BAD_REQUEST", "code is required")
	}
	err := h.svc.Bind(userID(c), body.Code)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAlreadyBound):
			return response.Fail(c, http.StatusConflict, "REFERRAL_ALREADY_BOUND", err.Error())
		case errors.Is(err, service.ErrBindSelf):
			return response.Fail(c, http.StatusBadRequest, "REFERRAL_BIND_SELF", err.Error())
		case errors.Is(err, service.ErrReferralCode):
			return response.Fail(c, http.StatusBadRequest, "REFERRAL_INVALID_CODE", err.Error())
		default:
			return response.Fail(c, http.StatusInternalServerError, "REFERRAL_BIND_FAILED", err.Error())
		}
	}
	return response.OK(c, map[string]bool{"bound": true})
}

func (h *ReferralHandler) Status(c echo.Context) error {
	data, err := h.svc.Status(userID(c))
	if err != nil {
		return response.Fail(c, http.StatusInternalServerError, "REFERRAL_STATUS_FAILED", err.Error())
	}
	return response.OK(c, data)
}
