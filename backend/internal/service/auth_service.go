package service

import (
	"errors"
	"time"

	"iot-rd-backend/internal/config"
	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo *repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{userRepo: userRepo, cfg: cfg}
}

type JWTClaims struct {
	UserID      string   `json:"userId"`
	Email       string   `json:"email"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

type LoginResponse struct {
	Token       string      `json:"token"`
	User        *model.User `json:"user"`
	Permissions []string    `json:"permissions"`
}

func (s *AuthService) Login(email, password string) (*LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if user.Status != "Active" {
		return nil, errors.New("user account is inactive or suspended")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Update last login
	now := time.Now()
	user.LastLoginAt = &now
	_ = s.userRepo.Update(user)

	// Extract permissions
	var permCodes []string
	if user.Role != nil {
		for _, p := range user.Role.Permissions {
			permCodes = append(permCodes, p.Code)
		}
	}

	// Generate JWT
	roleName := "Viewer"
	if user.Role != nil {
		roleName = user.Role.Name
	}

	claims := JWTClaims{
		UserID:      user.ID,
		Email:       user.Email,
		Role:        roleName,
		Permissions: permCodes,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.cfg.JWTExpires) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:       tokenString,
		User:        user,
		Permissions: permCodes,
	}, nil
}

func (s *AuthService) GetMe(userID string) (*model.User, []string, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, nil, err
	}

	var permCodes []string
	if user.Role != nil {
		for _, p := range user.Role.Permissions {
			permCodes = append(permCodes, p.Code)
		}
	}

	return user, permCodes, nil
}
