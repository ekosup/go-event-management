package models

import (
	"time"

	"gorm.io/gorm"
)

// Event represents an event created by a user
type Event struct {
	gorm.Model
	Name        string    `json:"name"`
	Description string    `json:"description"`
	DateTime    time.Time `json:"date_time"`
	UserID      uint      `json:"user_id"`
	User        User      `gorm:"foreignKey:UserID"`
	Guests      []Guest   `gorm:"foreignKey:EventID"`
}
