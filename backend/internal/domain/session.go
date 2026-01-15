package domain

import "time"

type TokenType string

const (
	TokenTypePlatform TokenType = "platform"
	TokenTypeSystem   TokenType = "system"
)

type Session struct {
	ID              int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	SessionToken    string     `json:"session_token" gorm:"unique;type:varchar(255)"`
	UserID          int64      `json:"user_id"`
	User            *User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	DeviceID        int64      `json:"device_id"`
	Device          *Device    `json:"device,omitempty" gorm:"foreignKey:DeviceID"`

	CurrentSystemID *int       `json:"current_system_id"`
	CurrentSystem   *System    `json:"current_system,omitempty" gorm:"foreignKey:CurrentSystemID"`

	OrganizationID  *int64        `json:"organization_id,omitempty"`
	Organization    *Organization `json:"organization,omitempty" gorm:"foreignKey:OrganizationID"`

	IPAddress        string     `json:"ip_address" gorm:"type:varchar(45)"`
	UserAgent        string     `json:"user_agent" gorm:"type:varchar(500)"`
	IsActive         *bool      `json:"is_active" gorm:"default:true"`
	ExpiresAt        time.Time  `json:"expires_at"`
	LastActivity     *time.Time `json:"last_activity"`
	LastSystemSwitch *time.Time `json:"last_system_switch"`
	LogoutAt         *time.Time `json:"logout_at"`
	LogoutReason     string     `json:"logout_reason" gorm:"type:varchar(100)"`
	ExtraFields
}

func (Session) TableName() string {
	return "sessions"
}

func (s *Session) IsInSystemContext() bool {
	return s.CurrentSystemID != nil
}
