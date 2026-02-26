package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"red_packet/backend/internal/http/response"
	"red_packet/backend/internal/service"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req service.LoginInput
	if err := c.Bind(&req); err != nil {
		return response.Fail(c, http.StatusBadRequest, "BAD_REQUEST", "invalid request")
	}
	token, user, err := h.svc.Login(req)
	if err != nil {
		return response.Fail(c, http.StatusInternalServerError, "AUTH_LOGIN_FAILED", err.Error())
	}
	return response.OK(c, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

func (h *AuthHandler) OTP(c echo.Context) error {
	type req struct {
		Account string `json:"account"`
	}
	var body req
	if err := c.Bind(&body); err != nil {
		return response.Fail(c, http.StatusBadRequest, "BAD_REQUEST", "invalid request")
	}
	return response.OK(c, map[string]string{"otp_masked": h.svc.OTP(body.Account)})
}
