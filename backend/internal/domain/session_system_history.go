package domain

import "time"

type SessionSystemHistory struct {
	ID         int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	SessionID  int64     `json:"session_id"`
	Session    *Session  `json:"session,omitempty" gorm:"foreignKey:SessionID"`
	SystemID   int       `json:"system_id"`
	System     *System   `json:"system,omitempty" gorm:"foreignKey:SystemID"`
	SwitchedAt time.Time `json:"switched_at"`
	IPAddress  string    `json:"ip_address" gorm:"type:varchar(45)"`
}

func (SessionSystemHistory) TableName() string {
	return "session_system_history"
}
