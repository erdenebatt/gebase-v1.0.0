package domain

type Module struct {
	ID            int            `json:"id" gorm:"primaryKey"`
	Code          string         `json:"code" gorm:"type:varchar(50)"`
	Name          string         `json:"name" gorm:"type:varchar(100)"`
	Description   string         `json:"description" gorm:"type:varchar(255)"`
	SystemID      *int           `json:"system_id"`
	System        *System        `json:"system,omitempty" gorm:"foreignKey:SystemID"`
	IsActive      *bool          `json:"is_active" gorm:"default:true"`
	ModuleActions []ModuleAction `json:"module_actions,omitempty" gorm:"foreignKey:ModuleID"`
	ExtraFields
}

func (Module) TableName() string {
	return "modules"
}

func (m *Module) IsPlatformModule() bool {
	return m.SystemID == nil
}
