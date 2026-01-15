package domain

import "time"

type DevicePlatform string

const (
	PlatformWeb            DevicePlatform = "web"
	PlatformIOS            DevicePlatform = "ios"
	PlatformAndroid        DevicePlatform = "android"
	PlatformTabletIOS      DevicePlatform = "tablet_ios"
	PlatformTabletAndroid  DevicePlatform = "tablet_android"
	PlatformWindowsDesktop DevicePlatform = "windows_desktop"
	PlatformMacDesktop     DevicePlatform = "mac_desktop"
	PlatformKiosk          DevicePlatform = "kiosk"
	PlatformPOSAndroid     DevicePlatform = "pos_android"
	PlatformPOSLinux       DevicePlatform = "pos_linux"
)

type Device struct {
	ID             int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	DeviceUID      string         `json:"device_uid" gorm:"unique;type:varchar(255)"`
	Name           string         `json:"name" gorm:"type:varchar(100)"`
	Platform       DevicePlatform `json:"platform" gorm:"type:varchar(50)"`
	OSVersion      string         `json:"os_version" gorm:"type:varchar(50)"`
	AppVersion     string         `json:"app_version" gorm:"type:varchar(50)"`
	PushToken      string         `json:"push_token" gorm:"type:varchar(500)"`
	OrganizationID *int64         `json:"organization_id,omitempty"`
	Organization   *Organization  `json:"organization,omitempty" gorm:"foreignKey:OrganizationID"`
	IsRegistered   *bool          `json:"is_registered" gorm:"default:false"`
	IsActive       *bool          `json:"is_active" gorm:"default:true"`
	RegisteredAt   *time.Time     `json:"registered_at"`
	LastHeartbeat  *time.Time     `json:"last_heartbeat"`
	ConfigJSON     string         `json:"config_json" gorm:"type:jsonb;default:'{}'"`
	ExtraFields
}

func (Device) TableName() string {
	return "devices"
}
