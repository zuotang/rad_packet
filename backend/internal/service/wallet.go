package service

import (
	"errors"

	"gorm.io/gorm"

	"red_packet/backend/internal/models"
)

type WalletView struct {
	Balance float64               `json:"balance"`
	Frozen  float64               `json:"frozen"`
	Ledgers []models.WalletLedger `json:"ledgers"`
}

type WalletService struct {
	db *gorm.DB
}

func NewWalletService(db *gorm.DB) *WalletService {
	return &WalletService{db: db}
}

func (s *WalletService) Get(userID uint, page, size int) (WalletView, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	var wallet models.Wallet
	err := s.db.Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			wallet = models.Wallet{UserID: userID}
			if e := s.db.Create(&wallet).Error; e != nil {
				return WalletView{}, e
			}
		} else {
			return WalletView{}, err
		}
	}

	var ledgers []models.WalletLedger
	if err := s.db.Where("user_id = ?", userID).
		Order("id DESC").
		Offset((page - 1) * size).
		Limit(size).
		Find(&ledgers).Error; err != nil {
		return WalletView{}, err
	}
	return WalletView{
		Balance: wallet.Balance,
		Frozen:  wallet.Frozen,
		Ledgers: ledgers,
	}, nil
}
