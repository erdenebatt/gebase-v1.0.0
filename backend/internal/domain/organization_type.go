package domain

type OrganizationType struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	Code        string `json:"code" gorm:"unique;type:varchar(50)"`
	Name        string `json:"name" gorm:"type:varchar(100)"`
	Description string `json:"description" gorm:"type:varchar(255)"`
	IsActive    *bool  `json:"is_active" gorm:"default:true"`
	ExtraFields
}

func (OrganizationType) TableName() string {
	return "organization_types"
}
