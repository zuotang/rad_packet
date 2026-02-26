package service

import (
	"gorm.io/gorm"

	"red_packet/backend/internal/models"
)

type BootstrapConfig struct {
	Tasks       []models.Task     `json:"tasks"`
	RewardTiers string            `json:"reward_tiers"`
	Configs     map[string]string `json:"configs"`
}

type ConfigService struct {
	db *gorm.DB
}

func NewConfigService(db *gorm.DB) *ConfigService {
	return &ConfigService{db: db}
}

func (s *ConfigService) Bootstrap() (BootstrapConfig, error) {
	var tasks []models.Task
	if err := s.db.Where("enabled = ?", true).Order("id ASC").Find(&tasks).Error; err != nil {
		return BootstrapConfig{}, err
	}
	var configs []models.AppConfig
	if err := s.db.Find(&configs).Error; err != nil {
		return BootstrapConfig{}, err
	}
	m := make(map[string]string, len(configs))
	rewardTiers := "[]"
	for _, c := range configs {
		m[c.Key] = c.Value
		if c.Key == "reward_tiers" {
			rewardTiers = c.Value
		}
	}
	return BootstrapConfig{Tasks: tasks, RewardTiers: rewardTiers, Configs: m}, nil
}

func (s *ConfigService) List() ([]models.AppConfig, error) {
	var items []models.AppConfig
	if err := s.db.Order("id ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (s *ConfigService) Upsert(key, value string) (models.AppConfig, error) {
	var cfg models.AppConfig
	if err := s.db.Where("`key` = ?", key).First(&cfg).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			cfg = models.AppConfig{Key: key, Value: value}
			if e := s.db.Create(&cfg).Error; e != nil {
				return models.AppConfig{}, e
			}
			return cfg, nil
		}
		return models.AppConfig{}, err
	}
	cfg.Value = value
	if err := s.db.Save(&cfg).Error; err != nil {
		return models.AppConfig{}, err
	}
	return cfg, nil
}
