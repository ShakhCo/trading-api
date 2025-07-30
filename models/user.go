package models

import "time"

type User struct {
	TelegramID int64     `gorm:"primaryKey"` // Primary key now
	FirstName  string    `gorm:"not null"`
	LastName   *string   // Optional
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}
