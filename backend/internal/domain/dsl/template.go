package dsl

import "gebase/internal/domain"

type TemplateType string

const (
	TemplateTypeEmail    TemplateType = "email"
	TemplateTypeSMS      TemplateType = "sms"
	TemplateTypePush     TemplateType = "push"
	TemplateTypeDocument TemplateType = "document"
	TemplateTypeReport   TemplateType = "report"
)

type Template struct {
	ID             int          `json:"id" gorm:"primaryKey"`
	Code           string       `json:"code" gorm:"unique;type:varchar(100)"`
	Name           string       `json:"name" gorm:"type:varchar(255)"`
	Description    string       `json:"description" gorm:"type:text"`
	TemplateType   TemplateType `json:"template_type" gorm:"type:varchar(50)"`
	Subject        string       `json:"subject" gorm:"type:varchar(500)"`
	Content        string       `json:"content" gorm:"type:text"`
	ContentHTML    string       `json:"content_html" gorm:"type:text"`
	Variables      string       `json:"variables" gorm:"type:jsonb"`
	OrganizationID *int64       `json:"organization_id,omitempty"`
	IsSystem       *bool        `json:"is_system" gorm:"default:false"`
	IsActive       *bool        `json:"is_active" gorm:"default:true"`
	domain.ExtraFields
}

func (Template) TableName() string {
	return "dsl_templates"
}
