package domain

type System struct {
	ID          int      `json:"id" gorm:"primaryKey"`
	Code        string   `json:"code" gorm:"unique;type:varchar(50)"` // dsl, gateway
	Name        string   `json:"name" gorm:"type:varchar(100)"`
	Description string   `json:"description" gorm:"type:varchar(255)"`
	IconURL     string   `json:"icon_url" gorm:"type:varchar(255)"`
	IconName    string   `json:"icon_name" gorm:"type:varchar(50)"`
	BaseURL     string   `json:"base_url" gorm:"type:varchar(255)"`
	Color       string   `json:"color" gorm:"type:varchar(20)"`
	IsActive    *bool    `json:"is_active" gorm:"default:true"`
	Sequence    int      `json:"sequence" gorm:"default:0"`
	Modules     []Module `json:"modules,omitempty" gorm:"foreignKey:SystemID"`
	Menus       []Menu   `json:"menus,omitempty" gorm:"foreignKey:SystemID"`
	ExtraFields
}

func (System) TableName() string {
	return "systems"
}
