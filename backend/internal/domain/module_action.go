package domain

type ModuleAction struct {
	ID       int     `json:"id" gorm:"primaryKey"`
	ModuleID int     `json:"module_id"`
	Module   *Module `json:"module,omitempty" gorm:"foreignKey:ModuleID"`
	ActionID int     `json:"action_id"`
	Action   *Action `json:"action,omitempty" gorm:"foreignKey:ActionID"`
	IsActive *bool   `json:"is_active" gorm:"default:true"`
	ExtraFields
}

func (ModuleAction) TableName() string {
	return "module_actions"
}
