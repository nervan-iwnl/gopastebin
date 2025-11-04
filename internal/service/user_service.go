package service

import (
	"context"
	"errors"

	"gopastebin/internal/domain"
	"gopastebin/internal/repository"
)

type UserService struct {
	users  repository.UserRepository
	pastes repository.PasteRepository
}

type UpdateUserDTO struct {
	Username *string `json:"username"`
	Avatar   *string `json:"avatar"`
}

func NewUserService(u repository.UserRepository, p repository.PasteRepository) *UserService {
	return &UserService{users: u, pastes: p}
}

func (s *UserService) GetPublicProfile(ctx context.Context, username string) (map[string]any, error) {
	user, err := s.users.FindByUsername(username)
	if err != nil {
		return nil, errors.New("not_found")
	}
	pastes, err := s.pastes.PublicOfUser(ctx, username, 50)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"user": map[string]any{
			"username": user.Username,
			"avatar":   user.Avatar,
		},
		"pastes": pastes,
	}, nil
}

func (s *UserService) GetMe(ctx context.Context, user any) (*domain.User, error) {
	uid := domain.ExtractUserID(user)
	return s.users.FindByID(uid)
}

func (s *UserService) UpdateMe(ctx context.Context, user any, dto UpdateUserDTO) (*domain.User, error) {
	uid := domain.ExtractUserID(user)
	u, err := s.users.FindByID(uid)
	if err != nil {
		return nil, err
	}
	if dto.Username != nil && *dto.Username != u.Username {
		if exists, _ := s.users.ExistsByUsername(*dto.Username); exists {
			return nil, errors.New("username_taken")
		}
		u.Username = *dto.Username
	}
	if dto.Avatar != nil {
		u.Avatar = *dto.Avatar
	}
	if err := s.users.Update(u); err != nil {
		return nil, err
	}
	return u, nil
}
