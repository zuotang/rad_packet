package service

import (
	"gorm.io/gorm"

	"red_packet/backend/internal/models"
)

type RiskService struct {
	db *gorm.DB
}

func NewRiskService(db *gorm.DB) *RiskService {
	return &RiskService{db: db}
}

func (s *RiskService) CheckWithdrawEligibility(userID uint) error {
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

func (s *RiskService) ListFlags(userID uint, page, size int) ([]models.RiskFlag, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	q := s.db.Model(&models.RiskFlag{})
	if userID > 0 {
		q = q.Where("user_id = ?", userID)
	}
	var items []models.RiskFlag
	if err := q.Order("id DESC").Offset((page - 1) * size).Limit(size).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (s *RiskService) AddFlag(userID uint, reason string, score int) (models.RiskFlag, error) {
	item := models.RiskFlag{UserID: userID, Reason: reason, Score: score}
	if err := s.db.Create(&item).Error; err != nil {
		return models.RiskFlag{}, err
	}
	return item, nil
}

func (s *RiskService) ListBlacklist(page, size int) ([]models.Blacklist, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	var items []models.Blacklist
	if err := s.db.Order("id DESC").Offset((page - 1) * size).Limit(size).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (s *RiskService) AddBlacklist(typ, value, note string) (models.Blacklist, error) {
	item := models.Blacklist{Type: typ, Value: value, Note: note}
	if err := s.db.Create(&item).Error; err != nil {
		return models.Blacklist{}, err
	}
	return item, nil
}
