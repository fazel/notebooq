package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;size:100" json:"username"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	Notes     []Note    `gorm:"constraint:OnDelete:CASCADE" json:"notes,omitempty"`
}

type Note struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Title     string    `gorm:"size:255" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
