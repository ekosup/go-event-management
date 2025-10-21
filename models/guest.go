package models

import (
	"gorm.io/gorm"
)

// Guest represents a guest attending an event
type Guest struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Attended bool   `json:"attended" gorm:"default:false"`
	EventID  uint   `json:"event_id"`
	QRCode   string `json:"qr_code"`
}
