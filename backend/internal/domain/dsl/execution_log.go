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
	ReferenceID    int             `json:"reference_id"`
	ReferenceCode  string          `json:"reference_code" gorm:"type:varchar(100)"`
	Status         ExecutionStatus `json:"status" gorm:"type:varchar(50)"`
	InputData      string          `json:"input_data" gorm:"type:jsonb"`
	OutputData     string          `json:"output_data" gorm:"type:jsonb"`
	ErrorMessage   string          `json:"error_message" gorm:"type:text"`
	ExecutionTime  int64           `json:"execution_time"`
	ExecutedAt     time.Time       `json:"executed_at"`
	ExecutedBy     *int64          `json:"executed_by"`
	OrganizationID *int64          `json:"organization_id,omitempty"`
	domain.ExtraFields
}

func (ExecutionLog) TableName() string {
	return "dsl_execution_logs"
}
