package service

import (
	"time"

	"gorm.io/gorm"

	"red_packet/backend/internal/models"
)

type AdminOpsService struct {
	db *gorm.DB
}

type AdminDashboard struct {
	TotalUsers         int64   `json:"total_users"`
	NewUsersToday      int64   `json:"new_users_today"`
	PendingWithdraws   int64   `json:"pending_withdraws"`
	PendingWithdrawAmt float64 `json:"pending_withdraw_amount"`
	RewardsPendingAmt  float64 `json:"rewards_pending_amount"`
	RewardsUnlockedAmt float64 `json:"rewards_unlocked_amount"`
}

func NewAdminOpsService(db *gorm.DB) *AdminOpsService {
	return &AdminOpsService{db: db}
}

func (s *AdminOpsService) Dashboard() (AdminDashboard, error) {
	var out AdminDashboard
	start := time.Now().Truncate(24 * time.Hour)
	end := start.Add(24 * time.Hour)

	if err := s.db.Model(&models.User{}).Count(&out.TotalUsers).Error; err != nil {
		return out, err
	}
	if err := s.db.Model(&models.User{}).Where("created_at >= ? AND created_at < ?", start, end).Count(&out.NewUsersToday).Error; err != nil {
		return out, err
	}
	if err := s.db.Model(&models.WithdrawRequest{}).Where("status = ?", "pending").Count(&out.PendingWithdraws).Error; err != nil {
		return out, err
	}
	if err := s.db.Model(&models.WithdrawRequest{}).Where("status = ?", "pending").Select("COALESCE(SUM(amount),0)").Scan(&out.PendingWithdrawAmt).Error; err != nil {
		return out, err
	}
	if err := s.db.Model(&models.Reward{}).Where("status = ?", "pending").Select("COALESCE(SUM(amount),0)").Scan(&out.RewardsPendingAmt).Error; err != nil {
		return out, err
	}
	if err := s.db.Model(&models.Reward{}).Where("status = ?", "unlocked").Select("COALESCE(SUM(amount),0)").Scan(&out.RewardsUnlockedAmt).Error; err != nil {
		return out, err
	}
	return out, nil
}
