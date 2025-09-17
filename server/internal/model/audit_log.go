package model

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement"`
	ActorID   *uuid.UUID `gorm:"type:uuid"`
	Action    string
	Resource  string
	Meta      string `gorm:"type:jsonb;default:'{}'"`
	CreatedAt time.Time
}
