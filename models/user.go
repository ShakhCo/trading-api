package models

import "time"

type User struct {
	TelegramID int64  `gorm:"primaryKey"`
	FirstName  string `gorm:"not null"`
	LastName   *string
	CreatedAt  time.Time   `gorm:"autoCreateTime"`
	Photos     []UserPhoto `gorm:"foreignKey:UserID;references:TelegramID;constraint:OnDelete:CASCADE"`
}

type UserPhoto struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    int64     `gorm:"not null"` // Foreign key to User.TelegramID
	FilePath  string    `gorm:"not null"` // Store file path or URL
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
