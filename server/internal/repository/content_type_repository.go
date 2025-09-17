// internal/repository/content_type_repository.go
package repository

import (
	"cms/server/internal/model"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ContentTypeRepository interface {
	Create(ctx context.Context, ct *model.ContentType) error
	List(ctx context.Context) ([]model.ContentType, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.ContentType, error)
	Update(ctx context.Context, id uuid.UUID, name, slug string) error
	Delete(ctx context.Context, id uuid.UUID) error
	AddField(ctx context.Context, field *model.ContentField) error
}

type contentTypeRepository struct {
	db *gorm.DB
}

func NewContentTypeRepository(db *gorm.DB) ContentTypeRepository {
	return &contentTypeRepository{db: db}
}

func (r *contentTypeRepository) Create(ctx context.Context, ct *model.ContentType) error {
	return r.db.WithContext(ctx).Create(ct).Error
}

func (r *contentTypeRepository) List(ctx context.Context) ([]model.ContentType, error) {
	var list []model.ContentType
	err := r.db.WithContext(ctx).
		Preload("Fields").
		Order("created_at DESC").
		Find(&list).Error
	return list, err
}

func (r *contentTypeRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.ContentType, error) {
	var ct model.ContentType
	err := r.db.WithContext(ctx).
		Preload("Fields").
		First(&ct, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &ct, nil
}

func (r *contentTypeRepository) Update(ctx context.Context, id uuid.UUID, name, slug string) error {
	return r.db.WithContext(ctx).Model(&model.ContentType{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"name": name,
			"slug": slug,
		}).Error
}

func (r *contentTypeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.ContentType{}, "id = ?", id).Error
}

func (r *contentTypeRepository) AddField(ctx context.Context, field *model.ContentField) error {
	return r.db.WithContext(ctx).Create(field).Error
}
