package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"gopastebin/config"
	"gopastebin/internal/domain"
	"gopastebin/internal/repository"
	"gopastebin/pkg/email"
	"gopastebin/pkg/jwtutil"
	"gopastebin/pkg/password"
)

type AuthService struct {
	users  repository.UserRepository
	config *config.Config
}

type RegisterDTO struct {
	Email    string `json:"email" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginDTO struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func NewAuthService(u repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{users: u, config: cfg}
}

func genVerifyCode() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func (a *AuthService) Register(dto RegisterDTO) (*domain.User, error) {
	if ok, _ := a.users.ExistsByEmail(dto.Email); ok {
		return nil, errors.New("email_taken")
	}
	if ok, _ := a.users.ExistsByUsername(dto.Username); ok {
		return nil, errors.New("username_taken")
	}

	hash, err := password.Hash(dto.Password)
	if err != nil {
		return nil, fmt.Errorf("hash_error: %w", err)
	}

	code := genVerifyCode()

	u := &domain.User{
		Email:           dto.Email,
		Username:        dto.Username,
		Password:        hash,
		EmailVerified:   false,
		EmailVerifyCode: code,
	}

	if err := a.users.Create(u); err != nil {
		return nil, err
	}

	// отправка письма
	go email.SendVerifyEmail(a.config, u.Email, code)

	return u, nil
}

func (a *AuthService) Login(dto LoginDTO) (string, string, *domain.User, error) {
	u, err := a.users.FindByEmailOrUsername(dto.Login)
	if err != nil {
		return "", "", nil, errors.New("invalid_credentials")
	}

	if !password.Check(dto.Password, u.Password) {
		return "", "", nil, errors.New("invalid_credentials")
	}

	// вот тут проверка подтверждения по почте
	if !u.EmailVerified {
		return "", "", nil, errors.New("email_not_verified")
	}

	access, err := jwtutil.CreateAccessToken(u.ID, u.Email, u.Username, a.config.JWTSecret)
	if err != nil {
		return "", "", nil, err
	}
	refresh, err := jwtutil.CreateRefreshToken(u.ID, u.Email, u.Username, a.config.JWTRefreshSecret)
	if err != nil {
		return "", "", nil, err
	}

	return access, refresh, u, nil
}

func (a *AuthService) VerifyEmail(ctx context.Context, emailStr, code string) error {
	u, err := a.users.FindByEmail(emailStr)
	if err != nil {
		return errors.New("not_found")
	}
	if u.EmailVerified {
		return nil
	}
	if u.EmailVerifyCode != code {
		return errors.New("invalid_code")
	}
	u.EmailVerified = true
	u.EmailVerifyCode = ""
	return a.users.Update(u)
}

func (a *AuthService) VerifyAccess(token string) (*domain.User, error) {
	claims, err := jwtutil.ParseToken(token, a.config.JWTSecret)
	if err != nil {
		return nil, err
	}
	return a.users.FindByID(claims.UserID)
}

func (a *AuthService) RefreshTokens(refresh string) (string, string, *domain.User, error) {
	claims, err := jwtutil.ParseToken(refresh, a.config.JWTRefreshSecret)
	if err != nil {
		return "", "", nil, err
	}
	u, err := a.users.FindByID(claims.UserID)
	if err != nil {
		return "", "", nil, err
	}
	access, err := jwtutil.CreateAccessToken(u.ID, u.Email, u.Username, a.config.JWTSecret)
	if err != nil {
		return "", "", nil, err
	}
	newRefresh, err := jwtutil.CreateRefreshToken(u.ID, u.Email, u.Username, a.config.JWTRefreshSecret)
	if err != nil {
		return "", "", nil, err
	}
	return access, newRefresh, u, nil
}
