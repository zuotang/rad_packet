package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"

	"red_packet/backend/internal/models"
)

type WithdrawService struct {
	db      *gorm.DB
	riskSvc *RiskService
}

func NewWithdrawService(db *gorm.DB, riskSvc *RiskService) *WithdrawService {
	return &WithdrawService{db: db, riskSvc: riskSvc}
}

func (s *WithdrawService) Apply(userID uint, amount float64) (models.WithdrawRequest, error) {
	var req models.WithdrawRequest
	if amount <= 0 {
		return req, ErrInvalidAmount
	}
	minAmount := s.loadWithdrawMin()
	if amount < minAmount {
		return req, ErrWithdrawBelowMin
	}
	if err := s.riskSvc.CheckWithdrawEligibility(userID); err != nil {
		return req, err
	}
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var wallet models.Wallet
		if err := tx.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
			return err
		}
		if wallet.Balance < amount {
			return ErrInsufficientFunds
		}
		wallet.Balance -= amount
		wallet.Frozen += amount
		if err := tx.Save(&wallet).Error; err != nil {
			return err
		}

		req = models.WithdrawRequest{
			UserID: userID,
			Amount: amount,
			Status: "pending",
		}
		if err := tx.Create(&req).Error; err != nil {
			return err
		}
		ledger := models.WalletLedger{
			UserID:  userID,
			Amount:  -amount,
			Type:    "withdraw_freeze",
			RefType: "withdraw_request",
			RefID:   fmt.Sprintf("%d", req.ID),
		}
		if err := tx.Create(&ledger).Error; err != nil {
			if isDuplicate(err) {
				return nil
			}
			return err
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, ErrInsufficientFunds) {
			return req, ErrInsufficientFunds
		}
		return req, err
	}
	return req, nil
}

func (s *WithdrawService) loadWithdrawMin() float64 {
	var cfg models.AppConfig
	if err := s.db.Where("`key` = ?", "withdraw_min").First(&cfg).Error; err == nil {
		if v, err := strconv.ParseFloat(cfg.Value, 64); err == nil && v > 0 {
			return v
		}
	}
	var rewardCfg models.AppConfig
	if err := s.db.Where("`key` = ?", "reward_tiers").First(&rewardCfg).Error; err == nil {
		if v := parseFirstTargetAmount(rewardCfg.Value); v > 0 {
			return v
		}
	}
	return 60
}

func parseFirstTargetAmount(raw string) float64 {
	type tier struct {
		Target float64 `json:"target"`
	}
	var tiers []tier
	if err := json.Unmarshal([]byte(raw), &tiers); err != nil {
		return 0
	}
	if len(tiers) == 0 {
		return 0
	}
	return tiers[0].Target
}

func (s *WithdrawService) UpdateStatus(requestID uint, status string, note string) (models.WithdrawRequest, error) {
	var req models.WithdrawRequest
	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&req, requestID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrWithdrawNotFound
			}
			return err
		}

		if !validWithdrawTransition(req.Status, status) {
			return ErrWithdrawState
		}
		oldStatus := req.Status
		req.Status = status
		req.Note = note
		if err := tx.Save(&req).Error; err != nil {
			return err
		}

		var wallet models.Wallet
		if err := tx.Where("user_id = ?", req.UserID).First(&wallet).Error; err != nil {
			return err
		}

		switch {
		case oldStatus == "pending" && status == "rejected":
			wallet.Balance += req.Amount
			wallet.Frozen -= req.Amount
			if wallet.Frozen < 0 {
				wallet.Frozen = 0
			}
			if err := tx.Save(&wallet).Error; err != nil {
				return err
			}
			ledger := models.WalletLedger{
				UserID:  req.UserID,
				Amount:  req.Amount,
				Type:    "withdraw_reject_refund",
				RefType: "withdraw_request_reject",
				RefID:   fmt.Sprintf("%d", req.ID),
			}
			if err := tx.Create(&ledger).Error; err != nil && !isDuplicate(err) {
				return err
			}
		case oldStatus == "approved" && status == "paid":
			wallet.Frozen -= req.Amount
			if wallet.Frozen < 0 {
				wallet.Frozen = 0
			}
			if err := tx.Save(&wallet).Error; err != nil {
				return err
			}
			ledger := models.WalletLedger{
				UserID:  req.UserID,
				Amount:  0,
				Type:    "withdraw_paid",
				RefType: "withdraw_request_paid",
				RefID:   fmt.Sprintf("%d", req.ID),
			}
			if err := tx.Create(&ledger).Error; err != nil && !isDuplicate(err) {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return models.WithdrawRequest{}, err
	}
	return req, nil
}

func validWithdrawTransition(from string, to string) bool {
	if from == "pending" && (to == "approved" || to == "rejected") {
		return true
	}
	if from == "approved" && to == "paid" {
		return true
	}
	return false
}

func (s *WithdrawService) ListByUser(userID uint, status string, page, size int) ([]models.WithdrawRequest, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	q := s.db.Where("user_id = ?", userID)
	if status != "" {
		q = q.Where("status = ?", status)
	}
	var records []models.WithdrawRequest
	if err := q.Order("id DESC").
		Offset((page - 1) * size).
		Limit(size).
		Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func (s *WithdrawService) ListAll(status string, page, size int) ([]models.WithdrawRequest, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	q := s.db.Model(&models.WithdrawRequest{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	var records []models.WithdrawRequest
	if err := q.Order("id DESC").
		Offset((page - 1) * size).
		Limit(size).
		Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}
