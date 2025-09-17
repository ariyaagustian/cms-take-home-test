package repository

import (
	"context"

	"cms/server/internal/model"

	"gorm.io/gorm"
)

type AuditRepository interface {
	Log(ctx context.Context, log *model.AuditLog) error
}

type auditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) AuditRepository {
	return &auditRepository{db: db}
}

func (r *auditRepository) Log(ctx context.Context, log *model.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}
