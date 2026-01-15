package domain

import (
	"sort"
	"time"
)

type User struct {
	ID           int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	SsoUserID    int64  `json:"sso_user_id" gorm:"uniqueIndex"`
	CivilID      int64  `json:"civil_id"`
	RegNo        string `json:"reg_no" gorm:"type:varchar(10);uniqueIndex"`
	FamilyName   string `json:"family_name" gorm:"type:varchar(80)"`
	LastName     string `json:"last_name" gorm:"type:varchar(150)"`
	FirstName    string `json:"first_name" gorm:"type:varchar(150)"`
	Gender       int    `json:"gender"`
	BirthDate    string `json:"birth_date" gorm:"type:varchar(10)"`
	PhoneNo      string `json:"phone_no" gorm:"type:varchar(8)"`
	Email        string `json:"email" gorm:"type:varchar(80);uniqueIndex"`
	PasswordHash string `json:"-" gorm:"type:varchar(255)"`
	AvatarURL    *string    `json:"avatar_url,omitempty" gorm:"type:varchar(500)"`
	IsActive     *bool      `json:"is_active" gorm:"default:true"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	LanguageCode string     `json:"language_code" gorm:"type:varchar(5);default:'mn'"`

	OrganizationID  *int64           `json:"organization_id,omitempty"`
	Organization    *Organization    `json:"organization,omitempty" gorm:"foreignKey:OrganizationID"`

	DefaultSystemID *int             `json:"default_system_id,omitempty"`
	DefaultSystem   *System          `json:"default_system,omitempty" gorm:"foreignKey:DefaultSystemID"`

	UserSystemRoles []UserSystemRole `json:"user_system_roles,omitempty" gorm:"foreignKey:UserID"`

	ExtraFields
}

func (User) TableName() string {
	return "users"
}

func (u *User) GetAvailableSystems() []System {
	systemMap := make(map[int]System)
	for _, usr := range u.UserSystemRoles {
		if usr.SystemID != nil && usr.System != nil && usr.IsActive != nil && *usr.IsActive {
			systemMap[*usr.SystemID] = *usr.System
		}
	}

	systems := make([]System, 0, len(systemMap))
	for _, s := range systemMap {
		systems = append(systems, s)
	}

	// Sort by sequence
	sort.Slice(systems, func(i, j int) bool {
		return systems[i].Sequence < systems[j].Sequence
	})

	return systems
}
