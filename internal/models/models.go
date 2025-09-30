package models

import "time"

type User struct {
	ID         uint `gorm:"primaryKey"`
	Username   string
	Password   string
	Email      string
	IsVerified bool
	VerifyCode string
}

type Note struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Title     string    `gorm:"size:255" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
