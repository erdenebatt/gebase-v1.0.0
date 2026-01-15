package dsl

import "gebase/internal/domain"

type VariableScope string

const (
	VariableScopeGlobal       VariableScope = "global"
	VariableScopeOrganization VariableScope = "organization"
	VariableScopeUser         VariableScope = "user"
	VariableScopeSession      VariableScope = "session"
)

type Variable struct {
	ID             int           `json:"id" gorm:"primaryKey"`
	Code           string        `json:"code" gorm:"type:varchar(100)"`
	Name           string        `json:"name" gorm:"type:varchar(255)"`
	Description    string        `json:"description" gorm:"type:text"`
	Scope          VariableScope `json:"scope" gorm:"type:varchar(50)"`
	ValueType      string        `json:"value_type" gorm:"type:varchar(50)"`
	Value          string        `json:"value" gorm:"type:text"`
	OrganizationID *int64        `json:"organization_id,omitempty"`
	UserID         *int64        `json:"user_id,omitempty"`
	IsEncrypted    *bool         `json:"is_encrypted" gorm:"default:false"`
	IsActive       *bool         `json:"is_active" gorm:"default:true"`
	domain.ExtraFields
}

func (Variable) TableName() string {
	return "dsl_variables"
}
