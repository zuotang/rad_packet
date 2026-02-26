package service

import (
	"gorm.io/gorm"

	"red_packet/backend/internal/config"
)

type Container struct {
	Auth     *AuthService
	Referral *ReferralService
	Reward   *RewardService
	Risk     *RiskService
	AdminOps *AdminOpsService
	Task     *TaskService
	Lottery  *LotteryService
	Wallet   *WalletService
	Withdraw *WithdrawService
	Config   *ConfigService
}

func NewContainer(db *gorm.DB, cfg config.Config) *Container {
	rewardSvc := NewRewardService(db)
	riskSvc := NewRiskService(db)
	referralSvc := NewReferralService(db, rewardSvc)
	lotterySvc := NewLotteryService(db, rewardSvc)
	return &Container{
		Auth:     NewAuthService(db, cfg),
		Referral: referralSvc,
		Reward:   rewardSvc,
		Risk:     riskSvc,
		AdminOps: NewAdminOpsService(db),
		Task:     NewTaskService(db, lotterySvc, referralSvc),
		Lottery:  lotterySvc,
		Wallet:   NewWalletService(db),
		Withdraw: NewWithdrawService(db, riskSvc),
		Config:   NewConfigService(db),
	}
}
