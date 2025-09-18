package repository

import (
	"context"

	"cms/server/internal/model"

	"gorm.io/gorm"
)

type MediaRepository interface {
	Save(ctx context.Context, asset *model.MediaAsset) error
	FindByID(ctx context.Context, id string) (*model.MediaAsset, error)
	List(ctx context.Context, limit, offset int) ([]model.MediaAsset, int64, error)
	Delete(ctx context.Context, id string) error
}

type mediaRepository struct {
	db *gorm.DB
}

func NewMediaRepository(db *gorm.DB) MediaRepository {
	return &mediaRepository{db: db}
}

func (r *mediaRepository) Save(ctx context.Context, asset *model.MediaAsset) error {
	return r.db.WithContext(ctx).Create(asset).Error
}

func (r *mediaRepository) FindByID(ctx context.Context, id string) (*model.MediaAsset, error) {
	var asset model.MediaAsset
	if err := r.db.WithContext(ctx).First(&asset, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *mediaRepository) List(ctx context.Context, limit, offset int) ([]model.MediaAsset, int64, error) {
	var items []model.MediaAsset
	var total int64

	if err := r.db.WithContext(ctx).
		Model(&model.MediaAsset{}).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *mediaRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Delete(&model.MediaAsset{}, "id = ?", id).Error
}
