package service

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"

	"red_packet/backend/internal/models"
)

type ReferralService struct {
	db        *gorm.DB
	rewardSvc *RewardService
}

type ReferralStatus struct {
	MyCode        string        `json:"my_code"`
	InviteCount   int64         `json:"invite_count"`
	DirectInvites []models.User `json:"direct_invites"`
}

func NewReferralService(db *gorm.DB, rewardSvc *RewardService) *ReferralService {
	return &ReferralService{db: db, rewardSvc: rewardSvc}
}

func (s *ReferralService) Bind(userID uint, code string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var inviterCode models.ReferralCode
		if err := tx.Where("code = ?", code).First(&inviterCode).Error; err != nil {
			return ErrReferralCode
		}
		if inviterCode.UserID == userID {
			return ErrBindSelf
		}

		edge := models.ReferralEdge{ParentUserID: inviterCode.UserID, ChildUserID: userID, Level: 1}
		if err := tx.Create(&edge).Error; err != nil {
			if isDuplicate(err) {
				return ErrAlreadyBound
			}
			return err
		}

		var parentEdge models.ReferralEdge
		if err := tx.Where("child_user_id = ? AND level = 1", inviterCode.UserID).First(&parentEdge).Error; err == nil {
			level2 := models.ReferralEdge{ParentUserID: parentEdge.ParentUserID, ChildUserID: userID, Level: 2}
			if err := tx.Create(&level2).Error; err != nil && !isDuplicate(err) {
				return err
			}
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return nil
	})
}

func (s *ReferralService) Status(userID uint) (ReferralStatus, error) {
	var status ReferralStatus
	var rc models.ReferralCode
	if err := s.db.Where("user_id = ?", userID).First(&rc).Error; err == nil {
		status.MyCode = rc.Code
	}

	var childIDs []uint
	if err := s.db.Model(&models.ReferralEdge{}).
		Where("parent_user_id = ? AND level = 1", userID).
		Pluck("child_user_id", &childIDs).Error; err != nil {
		return status, err
	}

	status.InviteCount = int64(len(childIDs))
	if len(childIDs) == 0 {
		return status, nil
	}

	if err := s.db.Where("id IN ?", childIDs).Order("id DESC").Limit(50).Find(&status.DirectInvites).Error; err != nil {
		return status, err
	}
	return status, nil
}

func (s *ReferralService) ProcessFirstValidAction(tx *gorm.DB, childUserID uint) error {
	var level1 models.ReferralEdge
	if err := tx.Where("child_user_id = ? AND level = 1", childUserID).First(&level1).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	if level1.IsValid {
		return nil
	}

	now := time.Now()
	level1.IsValid = true
	level1.ValidatedAt = &now
	if err := tx.Save(&level1).Error; err != nil {
		return err
	}

	amountL1, amountL2, err := s.loadInviteRewardRules(tx)
	if err != nil {
		return err
	}
	childIDRef := fmt.Sprintf("%d", childUserID)
	if amountL1 > 0 {
		if _, err := s.rewardSvc.GrantReward(tx, level1.ParentUserID, amountL1, "invite_valid_l1", childIDRef, "pending"); err != nil {
			return err
		}
	}

	var level2 models.ReferralEdge
	if err := tx.Where("child_user_id = ? AND level = 2", childUserID).First(&level2).Error; err == nil && amountL2 > 0 {
		if _, e := s.rewardSvc.GrantReward(tx, level2.ParentUserID, amountL2, "invite_valid_l2", childIDRef, "pending"); e != nil {
			return e
		}
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return nil
}

func (s *ReferralService) loadInviteRewardRules(tx *gorm.DB) (float64, float64, error) {
	var rows []models.AppConfig
	if err := tx.Where("`key` IN ?", []string{"invite_reward_l1", "invite_reward_l2"}).Find(&rows).Error; err != nil {
		return 0, 0, err
	}
	l1 := 3.0
	l2 := 1.0
	for _, row := range rows {
		if row.Key == "invite_reward_l1" {
			if v, err := strconv.ParseFloat(row.Value, 64); err == nil {
				l1 = v
			}
		}
		if row.Key == "invite_reward_l2" {
			if v, err := strconv.ParseFloat(row.Value, 64); err == nil {
				l2 = v
			}
		}
	}
	return l1, l2, nil
}
