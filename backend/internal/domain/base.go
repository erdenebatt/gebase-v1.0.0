package domain

import (
	"time"

	"gorm.io/gorm"
)

type ExtraFields struct {
	CreatedDate *time.Time     `json:"created_date,omitempty" gorm:"autoCreateTime"`
	UpdatedDate *time.Time     `json:"updated_date,omitempty" gorm:"autoUpdateTime"`
	CreatedBy   *int64         `json:"created_by,omitempty"`
	UpdatedBy   *int64         `json:"updated_by,omitempty"`
	DeletedDate gorm.DeletedAt `json:"deleted_date,omitempty" gorm:"index"`
	DeletedBy   *int64         `json:"deleted_by,omitempty"`
}

// Helper function for creating pointers
func Ptr[T any](v T) *T {
	return &v
}
