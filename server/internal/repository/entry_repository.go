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
	Delete(ctx context.Context, ctSlug string, id uuid.UUID, editorID *uuid.UUID) error
	Publish(ctx context.Context, ctSlug string, id uuid.UUID, t time.Time, editorID *uuid.UUID) error
	Rollback(ctx context.Context, ctSlug string, id uuid.UUID, version int, editorID *uuid.UUID) error
	ListPublished(ctx context.Context, slug string, limit, offset int, sort string) ([]model.Entry, int64, error)
	GetPublished(ctx context.Context, slug string, id uuid.UUID) (*model.Entry, error)
}

type entryRepository struct {
	db    *gorm.DB
	audit AuditRepository
}

func NewEntryRepository(db *gorm.DB, audit AuditRepository) EntryRepository {
	return &entryRepository{db: db, audit: audit}
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
	if err := r.db.WithContext(ctx).Create(&v).Error; err != nil {
		return err
	}

	// Audit log
	if r.audit != nil {
		_ = r.audit.Log(ctx, &model.AuditLog{
			ActorID:  editorID,
			Action:   "create_entry",
			Resource: "entry:" + e.ID.String(),
			Meta:     data,
		})
	}

	return nil
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
		if *status == "published" && e.PublishedAt == nil {
			now := time.Now()
			e.PublishedAt = &now
		}
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
	if err := r.db.WithContext(ctx).Create(&v).Error; err != nil {
		return err
	}

	// 👇 Tambahkan audit log
	if r.audit != nil {
		_ = r.audit.Log(ctx, &model.AuditLog{
			ActorID:  editorID,
			Action:   "update_entry",
			Resource: "entry:" + e.ID.String(),
			Meta:     data,
		})
	}

	return nil
}

func (r *entryRepository) Delete(ctx context.Context, slug string, id uuid.UUID, editorID *uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Entry{}).Error; err != nil {
		return err
	}

	// 👇 Tambahkan audit log
	if r.audit != nil {
		_ = r.audit.Log(ctx, &model.AuditLog{
			ActorID:  editorID,
			Action:   "delete_entry",
			Resource: "entry:" + id.String(),
			Meta:     json.RawMessage(`{"deleted": true}`),
		})
	}

	return nil
}

func (r *entryRepository) Publish(ctx context.Context, slug string, id uuid.UUID, t time.Time, editorID *uuid.UUID) error {
	entry, err := r.Get(ctx, slug, id)
	if err != nil {
		return err
	}
	status := "published"
	if err := r.Update(ctx, slug, id, entry.Data, &status, editorID); err != nil {
		return err
	}

	if r.audit != nil {
		_ = r.audit.Log(ctx, &model.AuditLog{
			ActorID:  editorID,
			Action:   "publish_entry",
			Resource: "entry:" + id.String(),
			Meta:     json.RawMessage(`{"slug": "` + slug + `"}`),
		})
	}

	return nil
}

func (r *entryRepository) Rollback(ctx context.Context, slug string, id uuid.UUID, version int, editorID *uuid.UUID) error {
	var v model.EntryVersion
	if err := r.db.WithContext(ctx).Where("entry_id = ? AND version = ?", id, version).First(&v).Error; err != nil {
		return err
	}

	if err := r.Update(ctx, slug, id, v.Data, nil, editorID); err != nil {
		return err
	}

	if r.audit != nil {
		meta, _ := json.Marshal(map[string]any{
			"slug":    slug,
			"version": version,
		})
		_ = r.audit.Log(ctx, &model.AuditLog{
			ActorID:  editorID,
			Action:   "rollback_entry",
			Resource: "entry:" + id.String(),
			Meta:     meta,
		})
	}

	return nil
}

func (r *entryRepository) ListPublished(ctx context.Context, slug string, limit, offset int, sort string) ([]model.Entry, int64, error) {
	ctID, err := r.findContentTypeID(slug)
	if err != nil {
		return nil, 0, err
	}

	var items []model.Entry
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Entry{}).
		Where("content_type_id = ? AND status = ? AND published_at <= ?", ctID, "published", time.Now())

	db.Count(&total)

	if sort == "-published_at" {
		db = db.Order("published_at desc")
	} else {
		db = db.Order("published_at asc")
	}

	if err := db.Limit(limit).Offset(offset).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *entryRepository) GetPublished(ctx context.Context, slug string, id uuid.UUID) (*model.Entry, error) {
	ctID, err := r.findContentTypeID(slug)
	if err != nil {
		return nil, err
	}
	var e model.Entry
	if err := r.db.WithContext(ctx).
		Where("content_type_id = ? AND id = ? AND status = ? AND published_at <= ?", ctID, id, "published", time.Now()).
		First(&e).Error; err != nil {
		return nil, err
	}
	return &e, nil
}
