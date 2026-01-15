package domain

type Organization struct {
	ID            int64             `json:"id,omitempty" gorm:"primaryKey;autoIncrement"`
	SsoOrgID      int64             `json:"sso_org_id" gorm:"uniqueIndex"`
	RegNo         string            `json:"reg_no,omitempty" gorm:"type:varchar(7);uniqueIndex"`
	Name          string            `json:"name,omitempty" gorm:"type:varchar(255)"`
	ShortName     string            `json:"short_name,omitempty" gorm:"type:varchar(255)"`
	TypeID        int               `json:"type_id"`
	Type          *OrganizationType `json:"type,omitempty" gorm:"foreignKey:TypeID"`
	PhoneNo       string            `json:"phone_no,omitempty" gorm:"type:varchar(8)"`
	Email         string            `json:"email,omitempty" gorm:"type:varchar(50)"`
	Longitude     float64           `json:"longitude,omitempty" gorm:"default:106.91758628931501"`
	Latitude      float64           `json:"latitude,omitempty" gorm:"default:47.918825014251915"`
	IsActive      *bool             `json:"is_active,omitempty" gorm:"default:true"`

	AimagID       int    `json:"aimag_id,omitempty"`
	SumID         int    `json:"sum_id,omitempty"`
	BagID         int    `json:"bag_id,omitempty"`
	AddressDetail string `json:"address_detail,omitempty" gorm:"type:varchar(255)"`
	AimagName     string `json:"aimag_name,omitempty" gorm:"type:varchar(255)"`
	SumName       string `json:"sum_name,omitempty" gorm:"type:varchar(255)"`
	BagName       string `json:"bag_name,omitempty" gorm:"type:varchar(255)"`
	CountryCode   string `json:"country_code,omitempty" gorm:"default:'MN'"`

	ParentID       *int64               `json:"parent_id"`
	Parent         *Organization        `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children       []Organization       `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Sequence       int                  `json:"sequence,omitempty"`
	EnabledSystems []OrganizationSystem `json:"enabled_systems,omitempty" gorm:"foreignKey:OrganizationID"`

	ExtraFields
}

func (Organization) TableName() string {
	return "organizations"
}
