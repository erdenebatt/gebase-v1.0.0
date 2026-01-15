package dsl

import (
	"time"

	"gebase/internal/domain"
)

type WorkflowStatus string

const (
	WorkflowStatusDraft     WorkflowStatus = "draft"
	WorkflowStatusPublished WorkflowStatus = "published"
	WorkflowStatusArchived  WorkflowStatus = "archived"
)

type Workflow struct {
	ID             int            `json:"id" gorm:"primaryKey"`
	Code           string         `json:"code" gorm:"unique;type:varchar(100)"`
	Name           string         `json:"name" gorm:"type:varchar(255)"`
	Description    string         `json:"description" gorm:"type:text"`
	SchemaID       *int           `json:"schema_id"`
	Schema         *Schema        `json:"schema,omitempty" gorm:"foreignKey:SchemaID"`
	Status         WorkflowStatus `json:"status" gorm:"type:varchar(50);default:'draft'"`
	TriggerType    string         `json:"trigger_type" gorm:"type:varchar(50)"`
	TriggerConfig  string         `json:"trigger_config" gorm:"type:jsonb"`
	OrganizationID *int64         `json:"organization_id,omitempty"`
	Version        int            `json:"version" gorm:"default:1"`
	IsActive       *bool          `json:"is_active" gorm:"default:true"`
	Steps          []WorkflowStep `json:"steps,omitempty" gorm:"foreignKey:WorkflowID"`
	domain.ExtraFields
}

func (Workflow) TableName() string {
	return "dsl_workflows"
}

type StepType string

const (
	StepTypeAction    StepType = "action"
	StepTypeCondition StepType = "condition"
	StepTypeLoop      StepType = "loop"
	StepTypeParallel  StepType = "parallel"
	StepTypeWait      StepType = "wait"
	StepTypeNotify    StepType = "notify"
)

type WorkflowStep struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	WorkflowID   int       `json:"workflow_id"`
	Workflow     *Workflow `json:"workflow,omitempty" gorm:"foreignKey:WorkflowID"`
	Code         string    `json:"code" gorm:"type:varchar(100)"`
	Name         string    `json:"name" gorm:"type:varchar(255)"`
	StepType     StepType  `json:"step_type" gorm:"type:varchar(50)"`
	Config       string    `json:"config" gorm:"type:jsonb"`
	NextStepCode string    `json:"next_step_code" gorm:"type:varchar(100)"`
	OnErrorCode  string    `json:"on_error_code" gorm:"type:varchar(100)"`
	Sequence     int       `json:"sequence" gorm:"default:0"`
	IsActive     *bool     `json:"is_active" gorm:"default:true"`
	domain.ExtraFields
}

func (WorkflowStep) TableName() string {
	return "dsl_workflow_steps"
}

type InstanceStatus string

const (
	InstanceStatusPending   InstanceStatus = "pending"
	InstanceStatusRunning   InstanceStatus = "running"
	InstanceStatusCompleted InstanceStatus = "completed"
	InstanceStatusFailed    InstanceStatus = "failed"
	InstanceStatusCancelled InstanceStatus = "cancelled"
)

type WorkflowInstance struct {
	ID           int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	WorkflowID   int            `json:"workflow_id"`
	Workflow     *Workflow      `json:"workflow,omitempty" gorm:"foreignKey:WorkflowID"`
	Status       InstanceStatus `json:"status" gorm:"type:varchar(50);default:'pending'"`
	InputData    string         `json:"input_data" gorm:"type:jsonb"`
	OutputData   string         `json:"output_data" gorm:"type:jsonb"`
	CurrentStep  string         `json:"current_step" gorm:"type:varchar(100)"`
	StartedAt    *time.Time     `json:"started_at"`
	CompletedAt  *time.Time     `json:"completed_at"`
	ErrorMessage string         `json:"error_message" gorm:"type:text"`
	TriggeredBy  *int64         `json:"triggered_by"`
	domain.ExtraFields
}

func (WorkflowInstance) TableName() string {
	return "dsl_workflow_instances"
}
