package repository

import (
	"context"

	"cms/server/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthRepository interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
	GetRoleByUserID(ctx context.Context, userID uuid.UUID) (string, error)
	AssignDefaultRole(ctx context.Context, userID uuid.UUID) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

// Cari user berdasarkan email
func (r *authRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Simpan user baru
func (r *authRepository) CreateUser(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// Ambil nama role dari user ID
func (r *authRepository) GetRoleByUserID(ctx context.Context, userID uuid.UUID) (string, error) {
	var userRole model.UserRole
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&userRole).Error
	if err != nil {
		return "", err
	}

	var role model.Role
	err = r.db.WithContext(ctx).Where("id = ?", userRole.RoleID).First(&role).Error
	if err != nil {
		return "", err
	}

	return role.Name, nil
}

// Default: assign role editor (id = 2) ke user baru
func (r *authRepository) AssignDefaultRole(ctx context.Context, userID uuid.UUID) error {
	role := model.UserRole{
		UserID: userID,
		RoleID: 2, // ⚠️ hardcoded role ID: "editor"
	}
	return r.db.WithContext(ctx).Create(&role).Error
}
