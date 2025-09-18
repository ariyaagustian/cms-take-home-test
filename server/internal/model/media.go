package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type MediaAsset struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Filename  string
	Mime      string
	SizeBytes int64
	URL       string
	Meta      json.RawMessage `gorm:"type:jsonb;default:'{}'"`
	CreatedBy *uuid.UUID      `gorm:"type:uuid"`
	CreatedAt time.Time
}
