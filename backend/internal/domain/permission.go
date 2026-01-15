package domain

type Permission struct {
	ID          int     `json:"id" gorm:"primaryKey"`
	Code        string  `json:"code" gorm:"unique;not null;type:varchar(255)"`
	Name        string  `json:"name" gorm:"type:varchar(255)"`
	Description string  `json:"description" gorm:"type:varchar(255)"`
	SystemID    *int    `json:"system_id"`
	System      *System `json:"system,omitempty" gorm:"foreignKey:SystemID"`
	ModuleID    int     `json:"module_id"`
	Module      *Module `json:"module,omitempty" gorm:"foreignKey:ModuleID"`
	ActionID    *int64  `json:"action_id"`
	Action      *Action `json:"action,omitempty" gorm:"foreignKey:ActionID"`
	IsActive    *bool   `json:"is_active" gorm:"default:true"`
	ExtraFields
}

func (Permission) TableName() string {
	return "permissions"
}

func (p *Permission) IsPlatformPermission() bool {
	return p.SystemID == nil
}
