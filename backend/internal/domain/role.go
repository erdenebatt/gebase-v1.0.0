package domain

type Role struct {
	ID          int              `json:"id" gorm:"primaryKey"`
	Code        string           `json:"code" gorm:"type:varchar(50)"`
	Name        string           `json:"name" gorm:"type:varchar(100)"`
	Description string           `json:"description" gorm:"type:varchar(255)"`
	SystemID    *int             `json:"system_id"`
	System      *System          `json:"system,omitempty" gorm:"foreignKey:SystemID"`
	IsSystem    *bool            `json:"is_system" gorm:"default:false"`
	IsActive    *bool            `json:"is_active" gorm:"default:true"`
	Permissions []RolePermission `json:"permissions,omitempty" gorm:"foreignKey:RoleID"`
	Menus       []RoleMenu       `json:"menus,omitempty" gorm:"foreignKey:RoleID"`
	ExtraFields
}

func (Role) TableName() string {
	return "roles"
}

func (r *Role) IsPlatformRole() bool {
	return r.SystemID == nil
}
