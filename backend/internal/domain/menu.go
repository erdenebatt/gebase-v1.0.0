package domain

type Menu struct {
	ID        int     `json:"id" gorm:"primaryKey"`
	Code      string  `json:"code" gorm:"type:varchar(50)"`
	Name      string  `json:"name" gorm:"type:varchar(100)"`
	SystemID  *int    `json:"system_id"`
	System    *System `json:"system,omitempty" gorm:"foreignKey:SystemID"`
	ParentID  *int    `json:"parent_id"`
	Parent    *Menu   `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children  []Menu  `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Path      string  `json:"path" gorm:"type:varchar(255)"`
	Icon      string  `json:"icon" gorm:"type:varchar(100)"`
	Component string  `json:"component" gorm:"type:varchar(255)"`
	Sequence  int     `json:"sequence" gorm:"default:0"`
	IsVisible *bool   `json:"is_visible" gorm:"default:true"`
	IsActive  *bool   `json:"is_active" gorm:"default:true"`
	ExtraFields
}

func (Menu) TableName() string {
	return "menus"
}
