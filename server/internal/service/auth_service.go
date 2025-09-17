package service

import (
	"context"
	"errors"
	"strings"

	"cms/server/internal/config"
	"cms/server/internal/model"
	"cms/server/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Authenticate(ctx context.Context, email, password string) (*model.AuthUser, error)
	Register(ctx context.Context, name, email, password string) (*model.AuthUser, error)
}

type authService struct {
	repo repository.AuthRepository
}

func NewAuthService(cfg config.Config, db *gorm.DB) AuthService {
	return &authService{
		repo: repository.NewAuthRepository(db),
	}
}

// ✅ Login
func (s *authService) Authenticate(ctx context.Context, email, password string) (*model.AuthUser, error) {
	email = strings.ToLower(email)

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	role, err := s.repo.GetRoleByUserID(ctx, user.ID)
	if err != nil {
		return nil, errors.New("failed to get user role")
	}

	return &model.AuthUser{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  role,
	}, nil
}

// ✅ Register
func (s *authService) Register(ctx context.Context, name, email, password string) (*model.AuthUser, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:         name,
		Email:        strings.ToLower(email),
		PasswordHash: string(hash),
	}

	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// Default role "editor" → ID = 2 (hardcoded)
	_ = s.repo.AssignDefaultRole(ctx, user.ID)

	return &model.AuthUser{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  "editor",
	}, nil
}
