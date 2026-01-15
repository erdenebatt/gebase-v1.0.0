package domain

type UserSystemRole struct {
	ID             int           `json:"id" gorm:"primaryKey"`
	UserID         int64         `json:"user_id"`
	User           *User         `json:"user,omitempty" gorm:"foreignKey:UserID"`
	SystemID       *int          `json:"system_id"`
	System         *System       `json:"system,omitempty" gorm:"foreignKey:SystemID"`
	RoleID         int           `json:"role_id"`
	Role           *Role         `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	OrganizationID *int64        `json:"organization_id,omitempty"`
	Organization   *Organization `json:"organization,omitempty" gorm:"foreignKey:OrganizationID"`
	IsActive       *bool         `json:"is_active" gorm:"default:true"`
	IsDefault      *bool         `json:"is_default" gorm:"default:false"`
	ExtraFields
}

func (UserSystemRole) TableName() string {
	return "user_system_roles"
}
