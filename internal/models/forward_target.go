package models

import (
	"time"
	
	"gorm.io/gorm"
)

type ForwardTarget struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"uniqueIndex;size:50" json:"name"`
	Email     string         `gorm:"size:100;not null" json:"email"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}