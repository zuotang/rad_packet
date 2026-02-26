package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"red_packet/backend/internal/http/response"
	"red_packet/backend/internal/models"
	"red_packet/backend/internal/service"
)

type AdminHandler struct {
	withdrawSvc *service.WithdrawService
	taskSvc     *service.TaskService
	configSvc   *service.ConfigService
	riskSvc     *service.RiskService
	opsSvc      *service.AdminOpsService
}

func NewAdminHandler(
	withdrawSvc *service.WithdrawService,
	taskSvc *service.TaskService,
	configSvc *service.ConfigService,
	riskSvc *service.RiskService,
	opsSvc *service.AdminOpsService,
) *AdminHandler {
	return &AdminHandler{
		withdrawSvc: withdrawSvc,
		taskSvc:     taskSvc,
		configSvc:   configSvc,
		riskSvc:     riskSvc,
		opsSvc:      opsSvc,
	}
}

func (h *AdminHandler) ReviewWithdraw(c echo.Context) error {
	type req struct {
		RequestID uint   `json:"request_id"`
		Status    string `json:"status"` // approved/rejected/paid
		Note      string `json:"note"`
	}
	var body req
	if err := c.Bind(&body); err != nil {
		return response.Fail(c, http.StatusBadRequest, "BAD_REQUEST", "invalid request")
	}
	data, err := h.withdrawSvc.UpdateStatus(body.RequestID, body.Status, body.Note)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrWithdrawNotFound):
			return response.Fail(c, http.StatusNotFound, "WITHDRAW_NOT_FOUND", err.Error())
		case errors.Is(err, service.ErrWithdrawState):
			return response.Fail(c, http.StatusBadRequest, "WITHDRAW_STATE_INVALID", err.Error())
		default:
			return response.Fail(c, http.StatusInternalServerError, "WITHDRAW_REVIEW_FAILED", err.Error())
		}
	}
	return response.OK(c, data)
}

func (h *AdminHandler) ListWithdraw(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	size, _ := strconv.Atoi(c.QueryParam("size"))
	status := strings.TrimSpace(c.QueryParam("status"))
	items, err := h.withdrawSvc.ListAll(status, page, size)
	if err != nil {
		return response.Fail(c, http.StatusInternalServerError, "ADMIN_WITHDRAW_LIST_FAILED", err.Error())
	}
	return response.OK(c, map[string]interface{}{"items": items})
}

func (h *AdminHandler) Dashboard(c echo.Context) error {
	data, err := h.opsSvc.Dashboard()
	if err != nil {
		return response.Fail(c, http.StatusInternalServerError, "ADMIN_DASHBOARD_FAILED", err.Error())
	}
	return response.OK(c, data)
}

func (h *AdminHandler) ListTasks(c echo.Context) error {
	items, err := h.taskSvc.ListAll()
	if err != nil {
		return response.Fail(c, http.StatusInternalServerError, "ADMIN_TASK_LIST_FAILED", err.Error())
	}
	return response.OK(c, map[string]interface{}{"items": items})
}

func (h *AdminHandler) SaveTask(c echo.Context) error {
	var in struct {
		ID           uint    `json:"id"`
		Type         string  `json:"type"`
		Name         string  `json:"name"`
		RewardRuleID string  `json:"reward_rule_id"`
		RewardAmount float64 `json:"reward_amount"`
		Enabled      bool    `json:"enabled"`
		CountryScope string  `json:"country_scope"`
	}
	if err := c.Bind(&in); err != nil {
		return response.Fail(c, http.StatusBadRequest, "BAD_REQUEST", "invalid request")
	}
	taskInput := models.Task{
		ID:           in.ID,
		Type:         in.Type,
		Name:         in.Name,
		RewardRuleID: in.RewardRuleID,
		RewardAmount: in.RewardAmount,
		Enabled:      in.Enabled,
		CountryScope: in.CountryScope,
	}
	task, err := h.taskSvc.SaveTask(taskInput)
	if err != nil {
		return response.Fail(c, http.StatusInternalServerError, "ADMIN_TASK_SAVE_FAILED", err.Error())
	}
	return response.OK(c, task)
}

func (h *AdminHandler) DeleteTask(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if id <= 0 {
		return response.Fail(c, http.StatusBadRequest, "BAD_REQUEST", "invalid id")
	}
	if err := h.taskSvc.DeleteTask(uint(id)); err != nil {
		return response.Fail(c, http.StatusInternalServerError, "ADMIN_TASK_DELETE_FAILED", err.Error())
	}
	return response.OK(c, map[string]bool{"deleted": true})
}

func (h *AdminHandler) ListConfigs(c echo.Context) error {
	items, err := h.configSvc.List()
	if err != nil {
		return response.Fail(c, http.StatusInternalServerError, "ADMIN_CONFIG_LIST_FAILED", err.Error())
	}
	return response.OK(c, map[string]interface{}{"items": items})
}

func (h *AdminHandler) UpsertConfig(c echo.Context) error {
	var in struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := c.Bind(&in); err != nil || strings.TrimSpace(in.Key) == "" {
		return response.Fail(c, http.StatusBadRequest, "BAD_REQUEST", "key is required")
	}
	item, err := h.configSvc.Upsert(strings.TrimSpace(in.Key), in.Value)
	if err != nil {
		return response.Fail(c, http.StatusInternalServerError, "ADMIN_CONFIG_SAVE_FAILED", err.Error())
	}
	return response.OK(c, item)
}

func (h *AdminHandler) ListRiskFlags(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	size, _ := strconv.Atoi(c.QueryParam("size"))
	userID, _ := strconv.Atoi(c.QueryParam("user_id"))
	items, err := h.riskSvc.ListFlags(uint(userID), page, size)
	if err != nil {
		return response.Fail(c, http.StatusInternalServerError, "ADMIN_RISK_FLAG_LIST_FAILED", err.Error())
	}
	return response.OK(c, map[string]interface{}{"items": items})
}

func (h *AdminHandler) AddRiskFlag(c echo.Context) error {
	var in struct {
		UserID uint   `json:"user_id"`
		Reason string `json:"reason"`
		Score  int    `json:"score"`
	}
	if err := c.Bind(&in); err != nil || in.UserID == 0 {
		return response.Fail(c, http.StatusBadRequest, "BAD_REQUEST", "invalid request")
	}
	if strings.TrimSpace(in.Reason) == "" {
		in.Reason = "manual admin flag"
	}
	item, err := h.riskSvc.AddFlag(in.UserID, in.Reason, in.Score)
	if err != nil {
		return response.Fail(c, http.StatusInternalServerError, "ADMIN_RISK_FLAG_ADD_FAILED", err.Error())
	}
	return response.OK(c, item)
}

func (h *AdminHandler) ListBlacklists(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	size, _ := strconv.Atoi(c.QueryParam("size"))
	items, err := h.riskSvc.ListBlacklist(page, size)
	if err != nil {
		return response.Fail(c, http.StatusInternalServerError, "ADMIN_BLACKLIST_LIST_FAILED", err.Error())
	}
	return response.OK(c, map[string]interface{}{"items": items})
}

func (h *AdminHandler) AddBlacklist(c echo.Context) error {
	var in struct {
		Type  string `json:"type"`
		Value string `json:"value"`
		Note  string `json:"note"`
	}
	if err := c.Bind(&in); err != nil || strings.TrimSpace(in.Type) == "" || strings.TrimSpace(in.Value) == "" {
		return response.Fail(c, http.StatusBadRequest, "BAD_REQUEST", "type/value required")
	}
	item, err := h.riskSvc.AddBlacklist(strings.TrimSpace(in.Type), strings.TrimSpace(in.Value), strings.TrimSpace(in.Note))
	if err != nil {
		return response.Fail(c, http.StatusInternalServerError, "ADMIN_BLACKLIST_ADD_FAILED", err.Error())
	}
	return response.OK(c, item)
}
