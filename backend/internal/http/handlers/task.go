package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"red_packet/backend/internal/http/response"
	"red_packet/backend/internal/service"
)

type TaskHandler struct {
	svc *service.TaskService
}

func NewTaskHandler(svc *service.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

func (h *TaskHandler) Claim(c echo.Context) error {
	var req service.ClaimInput
	if err := c.Bind(&req); err != nil {
		return response.Fail(c, http.StatusBadRequest, "BAD_REQUEST", "invalid request")
	}
	spinCount, err := h.svc.Claim(userID(c), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAlreadyClaimed):
			return response.Fail(c, http.StatusConflict, "TASK_ALREADY_CLAIMED", err.Error())
		case errors.Is(err, service.ErrTaskNotFound):
			return response.Fail(c, http.StatusNotFound, "TASK_NOT_FOUND", err.Error())
		default:
			return response.Fail(c, http.StatusInternalServerError, "TASK_CLAIM_FAILED", err.Error())
		}
	}
	return response.OK(c, map[string]interface{}{
		"claimed":        true,
		"spin_count":     spinCount,
	})
}

func (h *TaskHandler) List(c echo.Context) error {
	country := strings.TrimSpace(c.QueryParam("country"))
	items, err := h.svc.ListForUser(userID(c), country)
	if err != nil {
		return response.Fail(c, http.StatusInternalServerError, "TASK_LIST_FAILED", err.Error())
	}
	return response.OK(c, map[string]interface{}{"items": items})
}
