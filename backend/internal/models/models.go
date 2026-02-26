package models

import "time"

type User struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Phone      string    `gorm:"size:32;index" json:"phone"`
	Email      string    `gorm:"size:128;index" json:"email"`
	Country    string    `gorm:"size:16" json:"country"`
	Language   string    `gorm:"size:16" json:"language"`
	DeviceHash string    `gorm:"size:128;index" json:"device_hash"`
	CreatedAt  time.Time `json:"created_at"`
}

type ReferralCode struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"uniqueIndex"`
	Code      string `gorm:"size:32;uniqueIndex"`
	CreatedAt time.Time
}

type ReferralEdge struct {
	ID           uint `gorm:"primaryKey"`
	ParentUserID uint `gorm:"index"`
	ChildUserID  uint `gorm:"uniqueIndex:uniq_child_level"`
	Level        int  `gorm:"uniqueIndex:uniq_child_level"`
	IsValid      bool `gorm:"index;default:false"`
	ValidatedAt  *time.Time
	CreatedAt    time.Time
}

type Wallet struct {
	ID        uint    `gorm:"primaryKey"`
	UserID    uint    `gorm:"uniqueIndex"`
	Balance   float64 `gorm:"type:decimal(18,6);default:0"`
	Frozen    float64 `gorm:"type:decimal(18,6);default:0"`
	UpdatedAt time.Time
	CreatedAt time.Time
}

type WalletLedger struct {
	ID        uint    `gorm:"primaryKey"`
	UserID    uint    `gorm:"uniqueIndex:uniq_ledger"`
	Amount    float64 `gorm:"type:decimal(18,6)"`
	Type      string  `gorm:"size:32"`
	RefType   string  `gorm:"size:32;uniqueIndex:uniq_ledger"`
	RefID     string  `gorm:"size:64;uniqueIndex:uniq_ledger"`
	CreatedAt time.Time
}

type Reward struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       uint      `gorm:"index"`
	Status       string    `gorm:"size:16;index"` // pending/unlocked/expired
	Amount       float64   `gorm:"type:decimal(18,6)"`
	UnlockAmount float64   `gorm:"type:decimal(18,6)"`
	ExpireAt     time.Time `gorm:"index"`
	SourceType   string    `gorm:"size:32"`
	SourceID     string    `gorm:"size:64"`
	RewardType   string    `gorm:"size:32"`
	Description  string    `gorm:"size:255"`
	UnlockedAt   *time.Time
	CreatedAt    time.Time
}

type Task struct {
	ID           uint    `gorm:"primaryKey"`
	Type         string  `gorm:"size:32"`
	Name         string  `gorm:"size:64"`
	RewardRuleID string  `gorm:"size:64"`
	RewardAmount float64 `gorm:"type:decimal(18,6)"`
	Enabled      bool    `gorm:"index"`
	CountryScope string  `gorm:"size:255"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserTaskEvent struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"uniqueIndex:uniq_user_event"`
	TaskID    uint   `gorm:"index"`
	EventKey  string `gorm:"size:128;uniqueIndex:uniq_user_event"`
	MetaJSON  string `gorm:"type:text"`
	CreatedAt time.Time
}

type WithdrawRequest struct {
	ID        uint    `gorm:"primaryKey"`
	UserID    uint    `gorm:"index"`
	Amount    float64 `gorm:"type:decimal(18,6)"`
	Status    string  `gorm:"size:16;index"` // pending/approved/rejected/paid
	Note      string  `gorm:"size:255"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DeviceFingerprint struct {
	ID         uint   `gorm:"primaryKey"`
	UserID     uint   `gorm:"index"`
	DeviceHash string `gorm:"size:128;uniqueIndex"`
	FirstIP    string `gorm:"size:64"`
	LastIP     string `gorm:"size:64"`
	CreatedAt  time.Time
}

type RiskFlag struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"index"`
	Reason    string `gorm:"size:255"`
	Score     int
	CreatedAt time.Time
}

type Blacklist struct {
	ID        uint   `gorm:"primaryKey"`
	Type      string `gorm:"size:32;index"`
	Value     string `gorm:"size:255;index"`
	Note      string `gorm:"size:255"`
	CreatedAt time.Time
}

type AppConfig struct {
	ID        uint   `gorm:"primaryKey"`
	Key       string `gorm:"size:64;uniqueIndex"`
	Value     string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SpinChance struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"uniqueIndex"`
	Count     int  `gorm:"default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SpinRecord struct {
	ID           uint    `gorm:"primaryKey"`
	UserID       uint    `gorm:"index"`
	Amount       float64 `gorm:"type:decimal(18,6)"`
	PrizeType    string  `gorm:"size:32"`
	SegmentIndex int     `gorm:"index"`
	Status       string  `gorm:"size:16"` // win/lose
	CreatedAt    time.Time
}
