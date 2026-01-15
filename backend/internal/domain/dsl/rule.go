package dsl

import "gebase/internal/domain"

type RuleType string

const (
	RuleTypeValidation  RuleType = "validation"
	RuleTypeCalculation RuleType = "calculation"
	RuleTypeTrigger     RuleType = "trigger"
	RuleTypeConstraint  RuleType = "constraint"
)

type RuleTrigger string

const (
	TriggerBeforeInsert RuleTrigger = "before_insert"
	TriggerAfterInsert  RuleTrigger = "after_insert"
	TriggerBeforeUpdate RuleTrigger = "before_update"
	TriggerAfterUpdate  RuleTrigger = "after_update"
	TriggerBeforeDelete RuleTrigger = "before_delete"
	TriggerAfterDelete  RuleTrigger = "after_delete"
)

type Rule struct {
	ID           int         `json:"id" gorm:"primaryKey"`
	Code         string      `json:"code" gorm:"unique;type:varchar(100)"`
	Name         string      `json:"name" gorm:"type:varchar(255)"`
	Description  string      `json:"description" gorm:"type:text"`
	SchemaID     *int        `json:"schema_id"`
	Schema       *Schema     `json:"schema,omitempty" gorm:"foreignKey:SchemaID"`
	RuleType     RuleType    `json:"rule_type" gorm:"type:varchar(50)"`
	Trigger      RuleTrigger `json:"trigger" gorm:"type:varchar(50)"`
	Condition    string      `json:"condition" gorm:"type:text"`
	Expression   string      `json:"expression" gorm:"type:text"`
	ErrorMessage string      `json:"error_message" gorm:"type:varchar(500)"`
	Priority     int         `json:"priority" gorm:"default:0"`
	IsActive     *bool       `json:"is_active" gorm:"default:true"`
	domain.ExtraFields
}

func (Rule) TableName() string {
	return "dsl_rules"
}
