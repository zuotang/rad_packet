package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	"gorm.io/gorm"

	"red_packet/backend/internal/models"
)

type SpinResult struct {
	SpinID       uint    `json:"spin_id"`
	Amount       float64 `json:"amount"`
	PrizeType    string  `json:"prize_type"`
	SegmentIndex int     `json:"segment_index"`
	SpinCount    int     `json:"spin_count"`
}

type LotteryStatus struct {
	SpinCount  int     `json:"spin_count"`
	Target     float64 `json:"target"`
	Balance    float64 `json:"balance"`
	Pending    float64 `json:"pending"`
	Needed     float64 `json:"needed"`
	Unlockable float64 `json:"unlockable"`
}

type LotteryService struct {
	db        *gorm.DB
	rewardSvc *RewardService
}

func NewLotteryService(db *gorm.DB, rewardSvc *RewardService) *LotteryService {
	return &LotteryService{db: db, rewardSvc: rewardSvc}
}

func (s *LotteryService) GetStatus(userID uint) (LotteryStatus, error) {
	spinCount, err := s.GetSpinCount(userID)
	if err != nil {
		return LotteryStatus{}, err
	}
	target := s.loadWithdrawMin()
	balance, err := s.getBalance(userID)
	if err != nil {
		return LotteryStatus{}, err
	}
	summary, err := s.rewardSvc.Summary(userID)
	if err != nil {
		return LotteryStatus{}, err
	}
	needed := target - balance
	if needed < 0 {
		needed = 0
	}
	return LotteryStatus{
		SpinCount:  spinCount,
		Target:     target,
		Balance:    balance,
		Pending:    summary.Pending,
		Needed:     needed,
		Unlockable: summary.Pending,
	}, nil
}

func (s *LotteryService) Spin(userID uint) (SpinResult, error) {
	var result SpinResult
	target := s.loadWithdrawMin()
	return result, s.db.Transaction(func(tx *gorm.DB) error {
		chance, err := s.getOrCreateChance(tx, userID)
		if err != nil {
			return err
		}
		if chance.Count <= 0 {
			return ErrNoSpinChance
		}

		balance, err := s.getBalanceTx(tx, userID)
		if err != nil {
			return err
		}
		amount, prizeType := drawPrize(balance, target)
		segmentIndex := pickSegmentIndex(prizeType)

		record := models.SpinRecord{
			UserID:       userID,
			Amount:       amount,
			PrizeType:    prizeType,
			SegmentIndex: segmentIndex,
			Status:       "lose",
		}
		if amount > 0 {
			record.Status = "win"
		}
		if err := tx.Create(&record).Error; err != nil {
			return err
		}

		if amount > 0 {
			if _, err := s.rewardSvc.GrantReward(tx, userID, amount, "lottery_spin", fmt.Sprintf("%d", record.ID), "pending"); err != nil {
				return err
			}
		}

		chance.Count -= 1
		if chance.Count < 0 {
			chance.Count = 0
		}
		if err := tx.Save(&chance).Error; err != nil {
			return err
		}

		result = SpinResult{
			SpinID:       record.ID,
			Amount:       amount,
			PrizeType:    prizeType,
			SegmentIndex: segmentIndex,
			SpinCount:    chance.Count,
		}
		return nil
	})
}

func (s *LotteryService) AddChances(userID uint, count int) (int, error) {
	if count <= 0 {
		return 0, nil
	}
	return s.AddChancesTx(s.db, userID, count)
}

func (s *LotteryService) AddChancesTx(tx *gorm.DB, userID uint, count int) (int, error) {
	if count <= 0 {
		return 0, nil
	}
	var out int
	chance, err := s.getOrCreateChance(tx, userID)
	if err != nil {
		return 0, err
	}
	chance.Count += count
	if err := tx.Save(&chance).Error; err != nil {
		return 0, err
	}
	out = chance.Count
	return out, nil
}

func (s *LotteryService) GetSpinCount(userID uint) (int, error) {
	chance, err := s.getOrCreateChance(s.db, userID)
	if err != nil {
		return 0, err
	}
	return chance.Count, nil
}

func (s *LotteryService) ListRecords(userID uint, page, size int) ([]models.SpinRecord, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 20
	}
	var records []models.SpinRecord
	if err := s.db.Where("user_id = ?", userID).
		Order("id DESC").
		Offset((page - 1) * size).
		Limit(size).
		Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func (s *LotteryService) getOrCreateChance(tx *gorm.DB, userID uint) (models.SpinChance, error) {
	var chance models.SpinChance
	if err := tx.Where("user_id = ?", userID).First(&chance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			chance = models.SpinChance{UserID: userID, Count: 0}
			if err := tx.Create(&chance).Error; err != nil {
				return models.SpinChance{}, err
			}
			return chance, nil
		}
		return models.SpinChance{}, err
	}
	return chance, nil
}

func (s *LotteryService) getBalance(userID uint) (float64, error) {
	return s.getBalanceTx(s.db, userID)
}

func (s *LotteryService) getBalanceTx(tx *gorm.DB, userID uint) (float64, error) {
	var wallet models.Wallet
	if err := tx.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}
	return wallet.Balance, nil
}

func (s *LotteryService) loadWithdrawMin() float64 {
	var cfg models.AppConfig
	if err := s.db.Where("`key` = ?", "withdraw_min").First(&cfg).Error; err == nil {
		if v, err := strconv.ParseFloat(cfg.Value, 64); err == nil && v > 0 {
			return v
		}
	}
	var rewardCfg models.AppConfig
	if err := s.db.Where("`key` = ?", "reward_tiers").First(&rewardCfg).Error; err == nil {
		if v := parseFirstTarget(rewardCfg.Value); v > 0 {
			return v
		}
	}
	return 60
}

func parseFirstTarget(raw string) float64 {
	type tier struct {
		Target float64 `json:"target"`
	}
	var tiers []tier
	if err := jsonUnmarshal(raw, &tiers); err != nil {
		return 0
	}
	if len(tiers) == 0 {
		return 0
	}
	return tiers[0].Target
}

func jsonUnmarshal(raw string, out interface{}) error {
	if raw == "" {
		return errors.New("empty json")
	}
	return json.Unmarshal([]byte(raw), out)
}

func drawPrize(balance, target float64) (float64, string) {
	needed := target - balance
	progress := 0.0
	if target > 0 {
		progress = balance / target
	}
	if needed < 0 {
		needed = 0
	}

	stage := "mid"
	if progress < 0.3 {
		stage = "early"
	}
	if progress >= 0.8 || needed <= 1 {
		stage = "near"
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	p := r.Float64()

	switch stage {
	case "early":
		switch {
		case p < 0.15:
			return 0, "thanks"
		case p < 0.5:
			return randRange(r, 0.01, 0.3), "small"
		case p < 0.95:
			return randRange(r, 0.5, 3), "mid"
		default:
			return randRange(r, 5, 10), "big"
		}
	case "near":
		switch {
		case p < 0.1:
			return 0, "thanks"
		case p < 0.9:
			return randRange(r, 0.01, 0.2), "small"
		default:
			return randRange(r, 0.5, 1), "mid"
		}
	default:
		switch {
		case p < 0.12:
			return 0, "thanks"
		case p < 0.72:
			return randRange(r, 0.01, 0.3), "small"
		case p < 0.95:
			return randRange(r, 0.5, 3), "mid"
		default:
			return randRange(r, 5, 10), "big"
		}
	}
}

func randRange(r *rand.Rand, min, max float64) float64 {
	if min >= max {
		return round2(min)
	}
	val := min + r.Float64()*(max-min)
	if val < 0.01 {
		val = 0.01
	}
	return round2(val)
}

func round2(val float64) float64 {
	return math.Round(val*100) / 100
}

func pickSegmentIndex(prizeType string) int {
	segments := []string{"mid", "small", "mid", "thanks", "mid", "small", "mid", "thanks", "mid", "small", "big", "thanks"}
	var candidates []int
	for i, t := range segments {
		if t == prizeType {
			candidates = append(candidates, i)
		}
	}
	if len(candidates) == 0 {
		return 0
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return candidates[r.Intn(len(candidates))]
}
