package domain

type RolePermission struct {
	ID           int         `json:"id" gorm:"primaryKey"`
	RoleID       int         `json:"role_id"`
	Role         *Role       `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	PermissionID int         `json:"permission_id"`
	Permission   *Permission `json:"permission,omitempty" gorm:"foreignKey:PermissionID"`
	ExtraFields
}

func (RolePermission) TableName() string {
	return "role_permissions"
}
