package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"red_packet/backend/internal/models"
)

type GrantResult struct {
	Granted       bool    `json:"granted"`
	AlreadyExists bool    `json:"already_exists"`
	RewardID      uint    `json:"reward_id,omitempty"`
	Balance       float64 `json:"balance"`
	Frozen        float64 `json:"frozen"`
}

type RewardSummary struct {
	Pending  float64 `json:"pending"`
	Unlocked float64 `json:"unlocked"`
	Expired  float64 `json:"expired"`
}

type RewardService struct {
	db *gorm.DB
}

func NewRewardService(db *gorm.DB) *RewardService {
	return &RewardService{db: db}
}

func (s *RewardService) GrantReward(tx *gorm.DB, userID uint, amount float64, refType, refID string, status string) (GrantResult, error) {
	w := tx
	if w == nil {
		w = s.db
	}
	result := GrantResult{}
	err := w.Transaction(func(t *gorm.DB) error {
		ledger := models.WalletLedger{
			UserID:  userID,
			Amount:  amount,
			Type:    "reward",
			RefType: refType,
			RefID:   refID,
		}
		if err := t.Create(&ledger).Error; err != nil {
			if isDuplicate(err) {
				result.AlreadyExists = true
				return nil
			}
			return err
		}

		reward := models.Reward{
			UserID:       userID,
			Status:       status,
			Amount:       amount,
			UnlockAmount: amount,
			ExpireAt:     time.Now().Add(30 * 24 * time.Hour),
			SourceType:   refType,
			SourceID:     refID,
			RewardType:   "task",
			Description:  "auto grant",
		}
		if err := t.Create(&reward).Error; err != nil {
			return err
		}
		result.RewardID = reward.ID

		var wallet models.Wallet
		if err := t.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				wallet = models.Wallet{UserID: userID}
				if err := t.Create(&wallet).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
		if status == "pending" {
			wallet.Frozen += amount
		} else {
			wallet.Balance += amount
		}
		if err := t.Save(&wallet).Error; err != nil {
			return err
		}
		result.Balance = wallet.Balance
		result.Frozen = wallet.Frozen
		result.Granted = true
		return nil
	})
	return result, err
}

func (s *RewardService) Summary(userID uint) (RewardSummary, error) {
	var summary RewardSummary
	query := func(status string) (float64, error) {
		var val float64
		row := s.db.Model(&models.Reward{}).
			Where("user_id = ? AND status = ?", userID, status).
			Select("COALESCE(SUM(amount),0)").Row()
		if err := row.Scan(&val); err != nil {
			return 0, err
		}
		return val, nil
	}

	var err error
	if summary.Pending, err = query("pending"); err != nil {
		return summary, err
	}
	if summary.Unlocked, err = query("unlocked"); err != nil {
		return summary, err
	}
	if summary.Expired, err = query("expired"); err != nil {
		return summary, err
	}
	return summary, nil
}

func (s *RewardService) UnlockPendingRewards(userID uint) (int64, RewardSummary, error) {
	if err := s.checkUnlockEligibility(userID); err != nil {
		return 0, RewardSummary{}, err
	}

	var unlockedCount int64
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var rewards []models.Reward
		if err := tx.Where("user_id = ? AND status = ?", userID, "pending").Find(&rewards).Error; err != nil {
			return err
		}
		if len(rewards) == 0 {
			return nil
		}

		var total float64
		now := time.Now()
		for i := range rewards {
			total += rewards[i].Amount
			rewards[i].Status = "unlocked"
			rewards[i].UnlockedAt = &now
			if err := tx.Save(&rewards[i]).Error; err != nil {
				return err
			}
			ledger := models.WalletLedger{
				UserID:  userID,
				Amount:  rewards[i].Amount,
				Type:    "reward_unlock",
				RefType: "reward_unlock",
				RefID:   fmt.Sprintf("%d", rewards[i].ID),
			}
			if err := tx.Create(&ledger).Error; err != nil && !isDuplicate(err) {
				return err
			}
		}

		var wallet models.Wallet
		if err := tx.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
			return err
		}
		wallet.Frozen -= total
		if wallet.Frozen < 0 {
			wallet.Frozen = 0
		}
		wallet.Balance += total
		if err := tx.Save(&wallet).Error; err != nil {
			return err
		}
		unlockedCount = int64(len(rewards))
		return nil
	})
	if err != nil {
		return 0, RewardSummary{}, err
	}
	summary, err := s.Summary(userID)
	return unlockedCount, summary, err
}

func (s *RewardService) checkUnlockEligibility(userID uint) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return err
	}
	if user.DeviceHash != "" {
		var dupCount int64
		if err := s.db.Model(&models.User{}).Where("device_hash = ?", user.DeviceHash).Count(&dupCount).Error; err != nil {
			return err
		}
		if dupCount > 3 {
			return ErrRiskCheckFailed
		}
	}
	var riskScore int64
	if err := s.db.Model(&models.RiskFlag{}).Where("user_id = ?", userID).Select("COALESCE(SUM(score),0)").Scan(&riskScore).Error; err != nil {
		return err
	}
	if riskScore >= 100 {
		return ErrRiskCheckFailed
	}
	return nil
}

func (s *RewardService) ListByUser(userID uint, status string, page, size int) ([]models.Reward, error) {
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
	var records []models.Reward
	if err := q.Order("id DESC").
		Offset((page - 1) * size).
		Limit(size).
		Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func isDuplicate(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(strings.ToLower(err.Error()), "duplicate")
}
