package repository

import (
	"context"

	"cms/server/internal/model"

	"gorm.io/gorm"
)

type RoleRepository interface {
	List(ctx context.Context) ([]model.Role, error)
	Create(ctx context.Context, name string) error
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db}
}

func (r *roleRepository) List(ctx context.Context) ([]model.Role, error) {
	var roles []model.Role
	if err := r.db.WithContext(ctx).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *roleRepository) Create(ctx context.Context, name string) error {
	role := model.Role{Name: name}
	return r.db.WithContext(ctx).Create(&role).Error
}
