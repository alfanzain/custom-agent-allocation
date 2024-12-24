package models

import (
	"time"

	"gorm.io/gorm"
)

type Agent struct {
	ID          uint           `gorm:"primaryKey"`
	Name        string         `gorm:"type:varchar(255);not null"`
	CurrentLoad int            `gorm:"not null;default:0"`
	MaxLoad     int            `gorm:"not null;default:1"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
