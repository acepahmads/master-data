package repository

import (
	"time"

	"iot-rd-backend/internal/model"

	"gorm.io/gorm"
)

type APITokenRepository struct {
	db *gorm.DB
}

func NewAPITokenRepository(db *gorm.DB) *APITokenRepository {
	return &APITokenRepository{db: db}
}

func (r *APITokenRepository) Create(token *model.APIToken) error {
	return r.db.Create(token).Error
}

func (r *APITokenRepository) GetAll() ([]model.APIToken, error) {
	var tokens []model.APIToken
	err := r.db.Order("created_at desc").Find(&tokens).Error
	return tokens, err
}

func (r *APITokenRepository) GetByTokenHash(hash string) (*model.APIToken, error) {
	var token model.APIToken
	err := r.db.Where("token_hash = ? AND status = ?", hash, "Active").First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *APITokenRepository) UpdateLastUsed(id string) {
	now := time.Now()
	_ = r.db.Model(&model.APIToken{}).Where("id = ?", id).Update("last_used_at", &now).Error
}

func (r *APITokenRepository) Revoke(id string) error {
	return r.db.Model(&model.APIToken{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     "Revoked",
		"updated_at": time.Now(),
	}).Error
}
