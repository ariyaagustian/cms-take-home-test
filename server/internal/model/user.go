package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name         string
	Email        string `gorm:"uniqueIndex"`
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Role struct {
	ID   int    `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"uniqueIndex"`
}

type UserRole struct {
	UserID uuid.UUID `gorm:"type:uuid"`
	RoleID int
}
