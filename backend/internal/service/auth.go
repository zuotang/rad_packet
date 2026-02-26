package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"red_packet/backend/internal/config"
	"red_packet/backend/internal/models"
)

type LoginInput struct {
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	DeviceHash string `json:"device_hash"`
	Country    string `json:"country"`
	Language   string `json:"language"`
}

type AuthService struct {
	db  *gorm.DB
	cfg config.Config
}

func NewAuthService(db *gorm.DB, cfg config.Config) *AuthService {
	return &AuthService{db: db, cfg: cfg}
}

func (s *AuthService) Login(in LoginInput) (string, models.User, error) {
	var user models.User
	q := s.db.Model(&models.User{})
	if in.Phone != "" {
		q = q.Where("phone = ?", in.Phone)
	} else if in.Email != "" {
		q = q.Where("email = ?", in.Email)
	}

	err := q.First(&user).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return "", user, err
		}
		user = models.User{
			Phone:      in.Phone,
			Email:      in.Email,
			DeviceHash: in.DeviceHash,
			Country:    in.Country,
			Language:   in.Language,
		}
		if err := s.db.Create(&user).Error; err != nil {
			return "", user, err
		}
		rc := models.ReferralCode{
			UserID: user.ID,
			Code:   fmt.Sprintf("U%06d", user.ID),
		}
		if err := s.db.Create(&rc).Error; err != nil {
			return "", user, err
		}
		wallet := models.Wallet{UserID: user.ID}
	if err := s.db.Create(&wallet).Error; err != nil {
		return "", user, err
	}
	}

	var chance models.SpinChance
	if err := s.db.Where("user_id = ?", user.ID).First(&chance).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			chance = models.SpinChance{UserID: user.ID, Count: 100}
			if err := s.db.Create(&chance).Error; err != nil {
				return "", user, err
			}
		} else {
			return "", user, err
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": user.ID,
		"exp": time.Now().Add(time.Duration(s.cfg.JWT.TTLHours) * time.Hour).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return "", user, err
	}
	return tokenStr, user, nil
}

func (s *AuthService) OTP(_ string) string {
	return strings.Repeat("*", 6)
}
