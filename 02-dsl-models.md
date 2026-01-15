# DSL Domain Models

All DSL models are located in `internal/domain/dsl/` directory.

## Schema Model (dsl/schema.go)

```go
package dsl

import "gebase/internal/domain"

type Schema struct {
    ID             int       `json:"id" gorm:"primaryKey"`
    Code           string    `json:"code" gorm:"unique;type:varchar(100)"`
    Name           string    `json:"name" gorm:"type:varchar(255)"`
    Description    string    `json:"description" gorm:"type:text"`
    TableName      string    `json:"table_name" gorm:"type:varchar(100)"`
    OrganizationID *int64    `json:"organization_id,omitempty"`
    IsSystem       *bool     `json:"is_system" gorm:"default:false"`
    IsActive       *bool     `json:"is_active" gorm:"default:true"`
    Version        int       `json:"version" gorm:"default:1"`
    Fields         []Field   `json:"fields,omitempty" gorm:"foreignKey:SchemaID"`
    domain.ExtraFields
}

func (Schema) TableName() string {
    return "dsl_schemas"
}
```

## Field Model (dsl/field.go)

```go
package dsl

import "gebase/internal/domain"

type FieldType string

const (
    FieldTypeString    FieldType = "string"
    FieldTypeText      FieldType = "text"
    FieldTypeInteger   FieldType = "integer"
    FieldTypeDecimal   FieldType = "decimal"
    FieldTypeBoolean   FieldType = "boolean"
    FieldTypeDate      FieldType = "date"
    FieldTypeDatetime  FieldType = "datetime"
    FieldTypeJSON      FieldType = "json"
    FieldTypeReference FieldType = "reference"
    FieldTypeEnum      FieldType = "enum"
    FieldTypeFile      FieldType = "file"
)

type Field struct {
    ID            int       `json:"id" gorm:"primaryKey"`
    SchemaID      int       `json:"schema_id"`
    Schema        *Schema   `json:"schema,omitempty" gorm:"foreignKey:SchemaID"`
    Code          string    `json:"code" gorm:"type:varchar(100)"`
    Name          string    `json:"name" gorm:"type:varchar(255)"`
    FieldType     FieldType `json:"field_type" gorm:"type:varchar(50)"`
    DefaultValue  string    `json:"default_value" gorm:"type:varchar(500)"`
    IsRequired    *bool     `json:"is_required" gorm:"default:false"`
    IsUnique      *bool     `json:"is_unique" gorm:"default:false"`
    IsIndexed     *bool     `json:"is_indexed" gorm:"default:false"`
    MinLength     *int      `json:"min_length"`
    MaxLength     *int      `json:"max_length"`
    MinValue      *float64  `json:"min_value"`
    MaxValue      *float64  `json:"max_value"`
    Pattern       string    `json:"pattern" gorm:"type:varchar(500)"`       // Regex validation
    RefSchemaID   *int      `json:"ref_schema_id"`                          // For reference type
    RefSchema     *Schema   `json:"ref_schema,omitempty" gorm:"foreignKey:RefSchemaID"`
    EnumValues    string    `json:"enum_values" gorm:"type:jsonb"`          // ["value1", "value2"]
    Sequence      int       `json:"sequence" gorm:"default:0"`
    IsActive      *bool     `json:"is_active" gorm:"default:true"`
    domain.ExtraFields
}

func (Field) TableName() string {
    return "dsl_fields"
}

// Add unique constraint: schema_id + code
```

## Rule Model (dsl/rule.go)

```go
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
    Condition    string      `json:"condition" gorm:"type:text"`    // DSL expression for when to apply
    Expression   string      `json:"expression" gorm:"type:text"`   // DSL action/formula to execute
    ErrorMessage string      `json:"error_message" gorm:"type:varchar(500)"`
    Priority     int         `json:"priority" gorm:"default:0"`
    IsActive     *bool       `json:"is_active" gorm:"default:true"`
    domain.ExtraFields
}

func (Rule) TableName() string {
    return "dsl_rules"
}
```

## Workflow Model (dsl/workflow.go)

```go
package dsl

import "gebase/internal/domain"

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
    TriggerType    string         `json:"trigger_type" gorm:"type:varchar(50)"` // manual, auto, scheduled
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
    Config       string    `json:"config" gorm:"type:jsonb"`           // Step-specific configuration
    NextStepCode string    `json:"next_step_code" gorm:"type:varchar(100)"`
    OnErrorCode  string    `json:"on_error_code" gorm:"type:varchar(100)"`
    Sequence     int       `json:"sequence" gorm:"default:0"`
    IsActive     *bool     `json:"is_active" gorm:"default:true"`
    domain.ExtraFields
}

func (WorkflowStep) TableName() string {
    return "dsl_workflow_steps"
}
```

## Workflow Instance Model (dsl/workflow_instance.go)

```go
package dsl

import (
    "time"
    "gebase/internal/domain"
)

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
    TriggeredBy  *int64         `json:"triggered_by"` // UserID
    domain.ExtraFields
}

func (WorkflowInstance) TableName() string {
    return "dsl_workflow_instances"
}
```

## Template Model (dsl/template.go)

```go
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
    Subject        string       `json:"subject" gorm:"type:varchar(500)"`     // For email
    Content        string       `json:"content" gorm:"type:text"`             // Template body with {{variables}}
    ContentHTML    string       `json:"content_html" gorm:"type:text"`
    Variables      string       `json:"variables" gorm:"type:jsonb"`          // Available variables metadata
    OrganizationID *int64       `json:"organization_id,omitempty"`
    IsSystem       *bool        `json:"is_system" gorm:"default:false"`
    IsActive       *bool        `json:"is_active" gorm:"default:true"`
    domain.ExtraFields
}

func (Template) TableName() string {
    return "dsl_templates"
}
```

## Function Model (dsl/function.go)

```go
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
    Parameters   string       `json:"parameters" gorm:"type:jsonb"`         // [{name, type, required, default}]
    ReturnType   string       `json:"return_type" gorm:"type:varchar(50)"`
    Body         string       `json:"body" gorm:"type:text"`                // Function implementation/expression
    Example      string       `json:"example" gorm:"type:text"`             // Usage example
    IsActive     *bool        `json:"is_active" gorm:"default:true"`
    domain.ExtraFields
}

func (Function) TableName() string {
    return "dsl_functions"
}
```

## Variable Model (dsl/variable.go)

```go
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
    ValueType      string        `json:"value_type" gorm:"type:varchar(50)"` // string, number, boolean, json
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

// Add unique constraint: code + scope + organization_id + user_id
```

## Execution Log Model (dsl/execution_log.go)

```go
package dsl

import (
    "time"
    "gebase/internal/domain"
)

type ExecutionType string

const (
    ExecutionTypeRule     ExecutionType = "rule"
    ExecutionTypeWorkflow ExecutionType = "workflow"
    ExecutionTypeFunction ExecutionType = "function"
)

type ExecutionStatus string

const (
    ExecutionStatusSuccess ExecutionStatus = "success"
    ExecutionStatusError   ExecutionStatus = "error"
    ExecutionStatusSkipped ExecutionStatus = "skipped"
)

type ExecutionLog struct {
    ID             int64           `json:"id" gorm:"primaryKey;autoIncrement"`
    ExecutionType  ExecutionType   `json:"execution_type" gorm:"type:varchar(50)"`
    ReferenceID    int             `json:"reference_id"`    // Rule ID, Workflow ID, or Function ID
    ReferenceCode  string          `json:"reference_code" gorm:"type:varchar(100)"`
    Status         ExecutionStatus `json:"status" gorm:"type:varchar(50)"`
    InputData      string          `json:"input_data" gorm:"type:jsonb"`
    OutputData     string          `json:"output_data" gorm:"type:jsonb"`
    ErrorMessage   string          `json:"error_message" gorm:"type:text"`
    ExecutionTime  int64           `json:"execution_time"` // in milliseconds
    ExecutedAt     time.Time       `json:"executed_at"`
    ExecutedBy     *int64          `json:"executed_by"` // UserID
    OrganizationID *int64          `json:"organization_id,omitempty"`
    domain.ExtraFields
}

func (ExecutionLog) TableName() string {
    return "dsl_execution_logs"
}
```

## DSL Modules Summary

| Module    | Description                          | Key Features                              |
|-----------|--------------------------------------|-------------------------------------------|
| schema    | Dynamic data structure definition    | Table generation, versioning              |
| field     | Field definitions for schemas        | Type validation, constraints              |
| rule      | Business rules engine                | Validation, calculation, triggers         |
| workflow  | Process automation                   | Multi-step, conditions, parallel          |
| template  | Content templates                    | Email, SMS, documents with variables      |
| function  | Custom functions                     | Built-in and user-defined                 |
| variable  | Configuration variables              | Global, org, user, session scopes         |
| executor  | Rule/workflow execution engine       | Real-time processing                      |
| log       | Execution audit trail                | Performance tracking, debugging           |
