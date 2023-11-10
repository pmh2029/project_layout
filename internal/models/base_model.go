package models

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type BaseModel struct {
	CreatedAt time.Time             `gorm:"column:created_at;not null;type:timestamp;default:current_timestamp" mapstructure:"created_at" json:"created_at"`
	UpdatedAt time.Time             `gorm:"column:updated_at;not null;type:timestamp;default:current_timestamp" mapstructure:"updated_at" json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at;type:integer;default:0;index" mapstructure:"deleted_at" json:"deleted_at"`
}
