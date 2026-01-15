package domain

import "time"

type OrganizationSystem struct {
	ID             int           `json:"id" gorm:"primaryKey"`
	OrganizationID int64         `json:"organization_id"`
	Organization   *Organization `json:"organization,omitempty" gorm:"foreignKey:OrganizationID"`
	SystemID       int           `json:"system_id"`
	System         *System       `json:"system,omitempty" gorm:"foreignKey:SystemID"`
	IsActive       *bool         `json:"is_active" gorm:"default:true"`
	ActivatedAt    *time.Time    `json:"activated_at"`
	ExpiresAt      *time.Time    `json:"expires_at"`
	MaxUsers       *int          `json:"max_users"`
	Config         string        `json:"config" gorm:"type:jsonb"`
	ExtraFields
}

func (OrganizationSystem) TableName() string {
	return "organization_systems"
}
