package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Entry struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ContentTypeID uuid.UUID `gorm:"type:uuid;index"`
	Slug          string
	Status        string
	Data          json.RawMessage `gorm:"type:jsonb;default:'{}'"`
	PublishedAt   *time.Time
	CreatedBy     *uuid.UUID `gorm:"type:uuid"`
	UpdatedBy     *uuid.UUID `gorm:"type:uuid"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type EntryVersion struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	EntryID   uuid.UUID `gorm:"type:uuid;index"`
	Version   int
	Data      json.RawMessage `gorm:"type:jsonb;default:'{}'"`
	EditorID  *uuid.UUID      `gorm:"type:uuid"`
	CreatedAt time.Time
}
