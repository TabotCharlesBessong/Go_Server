package types

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	Password     string    `gorm:"not null" json:"-"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	IsVerified   bool      `gorm:"default:false" json:"is_verified"`
	LastLogin    time.Time `json:"last_login"`
	ResetToken   string    `json:"-"`
	ResetExpires time.Time `json:"-"`
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
