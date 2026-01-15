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
	ID           int       `json:"id" gorm:"primaryKey"`
	SchemaID     int       `json:"schema_id"`
	Schema       *Schema   `json:"schema,omitempty" gorm:"foreignKey:SchemaID"`
	Code         string    `json:"code" gorm:"type:varchar(100)"`
	Name         string    `json:"name" gorm:"type:varchar(255)"`
	FieldType    FieldType `json:"field_type" gorm:"type:varchar(50)"`
	DefaultValue string    `json:"default_value" gorm:"type:varchar(500)"`
	IsRequired   *bool     `json:"is_required" gorm:"default:false"`
	IsUnique     *bool     `json:"is_unique" gorm:"default:false"`
	IsIndexed    *bool     `json:"is_indexed" gorm:"default:false"`
	MinLength    *int      `json:"min_length"`
	MaxLength    *int      `json:"max_length"`
	MinValue     *float64  `json:"min_value"`
	MaxValue     *float64  `json:"max_value"`
	Pattern      string    `json:"pattern" gorm:"type:varchar(500)"`
	RefSchemaID  *int      `json:"ref_schema_id"`
	RefSchema    *Schema   `json:"ref_schema,omitempty" gorm:"foreignKey:RefSchemaID"`
	EnumValues   string    `json:"enum_values" gorm:"type:jsonb"`
	Sequence     int       `json:"sequence" gorm:"default:0"`
	IsActive     *bool     `json:"is_active" gorm:"default:true"`
	domain.ExtraFields
}

func (Field) TableName() string {
	return "dsl_fields"
}
