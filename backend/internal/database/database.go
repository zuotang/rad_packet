package database

import (
	"errors"

	"red_packet/backend/internal/config"
	"red_packet/backend/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Init(cfg config.Config) (*gorm.DB, error) {
	if cfg.MySQL.DSN == "" {
		return nil, errors.New("mysql dsn is empty")
	}
	db, err := gorm.Open(mysql.Open(cfg.MySQL.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.ReferralCode{},
		&models.ReferralEdge{},
		&models.Wallet{},
		&models.WalletLedger{},
		&models.Reward{},
		&models.Task{},
		&models.UserTaskEvent{},
		&models.WithdrawRequest{},
		&models.DeviceFingerprint{},
		&models.RiskFlag{},
		&models.Blacklist{},
		&models.AppConfig{},
		&models.SpinChance{},
		&models.SpinRecord{},
	); err != nil {
		return nil, err
	}

	if err := seed(db); err != nil {
		return nil, err
	}
	return db, nil
}

func seed(db *gorm.DB) error {
	defaultTasks := []models.Task{
		{Type: "checkin", Name: "每日签到", RewardRuleID: "daily_checkin", RewardAmount: 1, Enabled: true, CountryScope: "*"},
		{Type: "share", Name: "分享活动页", RewardRuleID: "share_landing", RewardAmount: 2, Enabled: true, CountryScope: "*"},
	}
	for _, t := range defaultTasks {
		if err := db.Where("reward_rule_id = ?", t.RewardRuleID).FirstOrCreate(&models.Task{}, t).Error; err != nil {
			return err
		}
	}
	bootstrap := models.AppConfig{
		Key:   "reward_tiers",
		Value: `[{"level":1,"target":50,"bonus":8},{"level":2,"target":150,"bonus":20}]`,
	}
	if err := db.Where("`key` = ?", bootstrap.Key).FirstOrCreate(&models.AppConfig{}, bootstrap).Error; err != nil {
		return err
	}
	defaultConfigs := []models.AppConfig{
		{Key: "invite_reward_l1", Value: "3"},
		{Key: "invite_reward_l2", Value: "1"},
		{Key: "withdraw_min", Value: "60"},
	}
	for _, c := range defaultConfigs {
		if err := db.Where("`key` = ?", c.Key).FirstOrCreate(&models.AppConfig{}, c).Error; err != nil {
			return err
		}
	}
	return nil
}
