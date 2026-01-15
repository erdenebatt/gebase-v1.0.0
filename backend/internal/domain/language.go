package domain

type Language struct {
	ID         int    `json:"id" gorm:"primaryKey"`
	Code       string `json:"code" gorm:"unique;type:varchar(5)"`
	Name       string `json:"name" gorm:"type:varchar(50)"`
	NativeName string `json:"native_name" gorm:"type:varchar(50)"`
	IsActive   *bool  `json:"is_active" gorm:"default:true"`
	IsDefault  *bool  `json:"is_default" gorm:"default:false"`
	ExtraFields
}

func (Language) TableName() string {
	return "languages"
}

type Translation struct {
	ID           int    `json:"id" gorm:"primaryKey"`
	LanguageCode string `json:"language_code" gorm:"type:varchar(5);index"`
	Key          string `json:"key" gorm:"type:varchar(255);index"`
	Value        string `json:"value" gorm:"type:text"`
	Module       string `json:"module" gorm:"type:varchar(50)"`
	ExtraFields
}

func (Translation) TableName() string {
	return "translations"
}
