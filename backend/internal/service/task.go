package service

import (
	"errors"
	"math"
	"strings"

	"gorm.io/gorm"

	"red_packet/backend/internal/models"
)

type ClaimInput struct {
	TaskID   uint   `json:"task_id"`
	EventKey string `json:"event_key"`
	MetaJSON string `json:"meta_json"`
}

type TaskService struct {
	db          *gorm.DB
	lotterySvc  *LotteryService
	referralSvc *ReferralService
}

type TaskView struct {
	ID                uint    `json:"id"`
	Type              string  `json:"type"`
	Name              string  `json:"name"`
	RewardAmount      float64 `json:"reward_amount"`
	Enabled           bool    `json:"enabled"`
	CountryScope      string  `json:"country_scope"`
	Claimed           bool    `json:"claimed"`
	LastClaimEventKey string  `json:"last_claim_event_key,omitempty"`
}

func NewTaskService(db *gorm.DB, lotterySvc *LotteryService, referralSvc *ReferralService) *TaskService {
	return &TaskService{db: db, lotterySvc: lotterySvc, referralSvc: referralSvc}
}

func (s *TaskService) Claim(userID uint, in ClaimInput) (int, error) {
	if in.EventKey == "" {
		return 0, ErrAlreadyClaimed
	}

	var spinCount int
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var task models.Task
		if err := tx.Where("id = ? AND enabled = ?", in.TaskID, true).First(&task).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrTaskNotFound
			}
			return err
		}

		var existed int64
		if err := tx.Model(&models.UserTaskEvent{}).
			Where("user_id = ? AND task_id = ?", userID, in.TaskID).
			Count(&existed).Error; err != nil {
			return err
		}
		if existed > 0 {
			return ErrAlreadyClaimed
		}

		ev := models.UserTaskEvent{
			UserID:   userID,
			TaskID:   in.TaskID,
			EventKey: in.EventKey,
			MetaJSON: in.MetaJSON,
		}
		if err := tx.Create(&ev).Error; err != nil {
			if isDuplicate(err) {
				return ErrAlreadyClaimed
			}
			return err
		}

		spinCount = int(math.Round(task.RewardAmount))
		if _, err := s.lotterySvc.AddChancesTx(tx, userID, spinCount); err != nil {
			return err
		}
		return s.referralSvc.ProcessFirstValidAction(tx, userID)
	})
	if err != nil {
		return 0, err
	}
	return spinCount, nil
}

func (s *TaskService) ListForUser(userID uint, country string) ([]TaskView, error) {
	var tasks []models.Task
	if err := s.db.Where("enabled = ?", true).Order("id ASC").Find(&tasks).Error; err != nil {
		return nil, err
	}

	var events []models.UserTaskEvent
	if err := s.db.Where("user_id = ?", userID).Order("id DESC").Find(&events).Error; err != nil {
		return nil, err
	}
	lastByTask := map[uint]string{}
	for _, ev := range events {
		if _, exists := lastByTask[ev.TaskID]; !exists {
			lastByTask[ev.TaskID] = ev.EventKey
		}
	}

	out := make([]TaskView, 0, len(tasks))
	for _, t := range tasks {
		if country != "" && t.CountryScope != "" && t.CountryScope != "*" && t.CountryScope != country {
			continue
		}
		ev, claimed := lastByTask[t.ID]
		out = append(out, TaskView{
			ID:                t.ID,
			Type:              t.Type,
			Name:              t.Name,
			RewardAmount:      t.RewardAmount,
			Enabled:           t.Enabled,
			CountryScope:      t.CountryScope,
			Claimed:           claimed,
			LastClaimEventKey: ev,
		})
	}
	return out, nil
}

func (s *TaskService) ListAll() ([]models.Task, error) {
	var tasks []models.Task
	if err := s.db.Order("id ASC").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *TaskService) SaveTask(in models.Task) (models.Task, error) {
	in.Type = strings.TrimSpace(in.Type)
	in.Name = strings.TrimSpace(in.Name)
	if in.Type == "" {
		in.Type = "custom"
	}
	if in.Name == "" {
		in.Name = "未命名任务"
	}
	if in.RewardRuleID == "" {
		in.RewardRuleID = in.Type
	}
	if in.CountryScope == "" {
		in.CountryScope = "*"
	}
	if in.ID == 0 {
		if err := s.db.Create(&in).Error; err != nil {
			return models.Task{}, err
		}
		return in, nil
	}
	if err := s.db.Model(&models.Task{}).Where("id = ?", in.ID).Updates(map[string]interface{}{
		"type":           in.Type,
		"name":           in.Name,
		"reward_rule_id": in.RewardRuleID,
		"reward_amount":  in.RewardAmount,
		"enabled":        in.Enabled,
		"country_scope":  in.CountryScope,
	}).Error; err != nil {
		return models.Task{}, err
	}
	var out models.Task
	if err := s.db.First(&out, in.ID).Error; err != nil {
		return models.Task{}, err
	}
	return out, nil
}

func (s *TaskService) DeleteTask(id uint) error {
	return s.db.Delete(&models.Task{}, id).Error
}
