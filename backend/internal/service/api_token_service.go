package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/repository"

	"github.com/google/uuid"
)

type APITokenService struct {
	repo         *repository.APITokenRepository
	activityRepo *repository.ActivityRepository
}

func NewAPITokenService(repo *repository.APITokenRepository, activityRepo *repository.ActivityRepository) *APITokenService {
	return &APITokenService{
		repo:         repo,
		activityRepo: activityRepo,
	}
}

type CreateTokenResponse struct {
	TokenInfo *model.APIToken `json:"tokenInfo"`
	RawToken  string          `json:"rawToken"` // Full token returned only once
}

func (s *APITokenService) GenerateToken(name string, expiresInDays int, scope string, user string) (*CreateTokenResponse, error) {
	if name == "" {
		name = "External API Client"
	}
	if scope == "" {
		scope = "ALL_READ"
	}

	// Generate 24 bytes cryptographically secure random token -> 48 hex chars
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	rawToken := fmt.Sprintf("iot_live_%s", hex.EncodeToString(b))

	// Hash token for database storage
	h := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(h[:])

	// Token preview (e.g. iot_live_9a8f...4e1b)
	preview := fmt.Sprintf("%s...%s", rawToken[:12], rawToken[len(rawToken)-4:])

	var expiresAt *time.Time
	if expiresInDays > 0 {
		exp := time.Now().AddDate(0, 0, expiresInDays)
		expiresAt = &exp
	} // if 0 or negative -> nil (Unlimited / Never Expire)

	token := &model.APIToken{
		ID:           fmt.Sprintf("TOK-%s", uuid.New().String()[:8]),
		Name:         name,
		TokenHash:    tokenHash,
		TokenPreview: preview,
		Scope:        scope,
		ExpiresAt:    expiresAt,
		CreatedBy:    user,
		Status:       "Active",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.Create(token); err != nil {
		return nil, err
	}

	_ = s.activityRepo.Create(&model.ActivityLog{
		ID:          uuid.New().String(),
		UserID:      user,
		Module:      "system",
		Action:      "GENERATE_API_TOKEN",
		EntityType:  "API_TOKEN",
		EntityID:    token.ID,
		Description: fmt.Sprintf("Generated API Token '%s' (%s) with expiration %v", token.Name, preview, expiresAt),
		CreatedAt:   time.Now(),
	})

	return &CreateTokenResponse{
		TokenInfo: token,
		RawToken:  rawToken,
	}, nil
}

func (s *APITokenService) GetAll() ([]model.APIToken, error) {
	return s.repo.GetAll()
}

func (s *APITokenService) Revoke(id string, user string) error {
	_ = s.activityRepo.Create(&model.ActivityLog{
		ID:          uuid.New().String(),
		UserID:      user,
		Module:      "system",
		Action:      "REVOKE_API_TOKEN",
		EntityType:  "API_TOKEN",
		EntityID:    id,
		Description: fmt.Sprintf("Revoked API Token ID %s", id),
		CreatedAt:   time.Now(),
	})
	return s.repo.Revoke(id)
}

func (s *APITokenService) ValidateToken(rawToken string) (*model.APIToken, error) {
	h := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(h[:])

	token, err := s.repo.GetByTokenHash(tokenHash)
	if err != nil {
		return nil, fmt.Errorf("invalid or revoked API token")
	}

	if token.ExpiresAt != nil && token.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("API token has expired")
	}

	go s.repo.UpdateLastUsed(token.ID)
	return token, nil
}
