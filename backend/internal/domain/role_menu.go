package domain

type RoleMenu struct {
	ID     int   `json:"id" gorm:"primaryKey"`
	RoleID int   `json:"role_id"`
	Role   *Role `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	MenuID int   `json:"menu_id"`
	Menu   *Menu `json:"menu,omitempty" gorm:"foreignKey:MenuID"`
	ExtraFields
}

func (RoleMenu) TableName() string {
	return "role_menus"
}
