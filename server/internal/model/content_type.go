package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ContentType struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `gorm:"uniqueIndex"`
	Slug      string    `gorm:"uniqueIndex"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Fields    []ContentField `gorm:"foreignKey:ContentTypeID"`
}

type ContentField struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ContentTypeID uuid.UUID `gorm:"type:uuid;index"`
	Name          string
	Kind          string
	Options       json.RawMessage `gorm:"type:jsonb;default:'{}'"`
}
