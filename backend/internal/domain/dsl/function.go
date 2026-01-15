package dsl

import "gebase/internal/domain"

type FunctionType string

const (
	FunctionTypeBuiltIn FunctionType = "builtin"
	FunctionTypeCustom  FunctionType = "custom"
)

type Function struct {
	ID           int          `json:"id" gorm:"primaryKey"`
	Code         string       `json:"code" gorm:"unique;type:varchar(100)"`
	Name         string       `json:"name" gorm:"type:varchar(255)"`
	Description  string       `json:"description" gorm:"type:text"`
	FunctionType FunctionType `json:"function_type" gorm:"type:varchar(50)"`
	Parameters   string       `json:"parameters" gorm:"type:jsonb"`
	ReturnType   string       `json:"return_type" gorm:"type:varchar(50)"`
	Body         string       `json:"body" gorm:"type:text"`
	Example      string       `json:"example" gorm:"type:text"`
	IsActive     *bool        `json:"is_active" gorm:"default:true"`
	domain.ExtraFields
}

func (Function) TableName() string {
	return "dsl_functions"
}
