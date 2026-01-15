package dsl

import "gebase/internal/domain"

type Schema struct {
	ID             int     `json:"id" gorm:"primaryKey"`
	Code           string  `json:"code" gorm:"unique;type:varchar(100)"`
	Name           string  `json:"name" gorm:"type:varchar(255)"`
	Description    string  `json:"description" gorm:"type:text"`
	DBTableName    string  `json:"table_name" gorm:"column:table_name;type:varchar(100)"`
	OrganizationID *int64  `json:"organization_id,omitempty"`
	IsSystem       *bool   `json:"is_system" gorm:"default:false"`
	IsActive       *bool   `json:"is_active" gorm:"default:true"`
	Version        int     `json:"version" gorm:"default:1"`
	Fields         []Field `json:"fields,omitempty" gorm:"foreignKey:SchemaID"`
	domain.ExtraFields
}

func (Schema) TableName() string {
	return "dsl_schemas"
}
