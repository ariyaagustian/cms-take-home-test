// internal/repository/entry_repository.go
package repository

import (
	"context"
	"encoding/json"
	"time"

	"cms/server/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EntryRepository interface {
	Create(ctx context.Context, ctSlug string, e *model.Entry, data json.RawMessage, editorID *uuid.UUID) error
	List(ctx context.Context, ctSlug string, limit, offset int) ([]model.Entry, int64, error)
	Get(ctx context.Context, ctSlug string, id uuid.UUID) (*model.Entry, error)
	Update(ctx context.Context, ctSlug string, id uuid.UUID, data json.RawMessage, status *string, editorID *uuid.UUID) error
	Delete(ctx context.Context, ctSlug string, id uuid.UUID) error
	Publish(ctx context.Context, ctSlug string, id uuid.UUID, t time.Time, editorID *uuid.UUID) error
	Rollback(ctx context.Context, ctSlug string, id uuid.UUID, version int, editorID *uuid.UUID) error
}

type entryRepository struct{ db *gorm.DB }

func NewEntryRepository(db *gorm.DB) EntryRepository {
	return &entryRepository{db: db}
}

func (r *entryRepository) findContentTypeID(slug string) (uuid.UUID, error) {
	var ct model.ContentType
	if err := r.db.Where("slug = ?", slug).First(&ct).Error; err != nil {
		return uuid.UUID{}, err
	}
	return ct.ID, nil
}

func (r *entryRepository) Create(ctx context.Context, slug string, e *model.Entry, data json.RawMessage, editorID *uuid.UUID) error {
	ctID, err := r.findContentTypeID(slug)
	if err != nil {
		return err
	}
	e.ContentTypeID = ctID
	e.Data = data
	e.CreatedBy = editorID
	e.UpdatedBy = editorID
	if err := r.db.WithContext(ctx).Create(e).Error; err != nil {
		return err
	}

	v := model.EntryVersion{
		EntryID:  e.ID,
		Version:  1,
		Data:     data,
		EditorID: editorID,
	}
	return r.db.WithContext(ctx).Create(&v).Error
}

func (r *entryRepository) List(ctx context.Context, slug string, limit, offset int) ([]model.Entry, int64, error) {
	ctID, err := r.findContentTypeID(slug)
	if err != nil {
		return nil, 0, err
	}
	var items []model.Entry
	var total int64
	db := r.db.WithContext(ctx).Where("content_type_id = ?", ctID)
	db.Model(&model.Entry{}).Count(&total)
	if err := db.Order("created_at desc").Limit(limit).Offset(offset).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *entryRepository) Get(ctx context.Context, slug string, id uuid.UUID) (*model.Entry, error) {
	ctID, err := r.findContentTypeID(slug)
	if err != nil {
		return nil, err
	}
	var e model.Entry
	if err := r.db.WithContext(ctx).Where("content_type_id = ? AND id = ?", ctID, id).First(&e).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *entryRepository) Update(ctx context.Context, slug string, id uuid.UUID, data json.RawMessage, status *string, editorID *uuid.UUID) error {
	e, err := r.Get(ctx, slug, id)
	if err != nil {
		return err
	}

	if data != nil {
		e.Data = data
	}

	e.Data = data
	e.UpdatedBy = editorID
	e.UpdatedAt = time.Now()
	if status != nil {
		e.Status = *status
	}
	if err := r.db.WithContext(ctx).Save(e).Error; err != nil {
		return err
	}

	var latestVer model.EntryVersion
	if err := r.db.WithContext(ctx).Where("entry_id = ?", id).Order("version desc").First(&latestVer).Error; err != nil {
		return err
	}

	v := model.EntryVersion{
		EntryID:  id,
		Version:  latestVer.Version + 1,
		Data:     data,
		EditorID: editorID,
	}
	return r.db.WithContext(ctx).Create(&v).Error
}

func (r *entryRepository) Delete(ctx context.Context, slug string, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Entry{}).Error
}

func (r *entryRepository) Publish(ctx context.Context, slug string, id uuid.UUID, t time.Time, editorID *uuid.UUID) error {
	entry, err := r.Get(ctx, slug, id)
	if err != nil {
		return err
	}
	s := "published"
	return r.Update(ctx, slug, id, entry.Data, &s, editorID)
}

func (r *entryRepository) Rollback(ctx context.Context, slug string, id uuid.UUID, version int, editorID *uuid.UUID) error {
	var v model.EntryVersion
	if err := r.db.WithContext(ctx).Where("entry_id = ? AND version = ?", id, version).First(&v).Error; err != nil {
		return err
	}
	return r.Update(ctx, slug, id, v.Data, nil, editorID)
}
