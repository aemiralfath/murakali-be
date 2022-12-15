package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id" db:"id" binding:"omitempty"`
	RoleID    int       `json:"role_id" db:"role_id" binding:"omitempty"`
	Username  string    `json:"username" db:"username" binding:"omitempty"`
	Email     string    `json:"email" db:"email" binding:"omitempty"`
	PhoneNo   string    `json:"phone_no" db:"phone_no" binding:"omitempty"`
	FullName  string    `json:"fullname" db:"fullname" binding:"omitempty"`
	Password  string    `json:"password" db:"password" binding:"omitempty"`
	Gender    string    `json:"gender" db:"gender" binding:"omitempty"`
	PhotoURL  string    `json:"photo_url" db:"photo_url" binding:"omitempty"`
	IsSSO     bool      `json:"is_sso" db:"is_sso" binding:"omitempty"`
	IsVerify  bool      `json:"is_verify" db:"is_verify" binding:"omitempty"`
	BirthDate time.Time `json:"birth_date" db:"birth_date" binding:"omitempty"`
	CreatedAt time.Time `json:"created_at" db:"created_at" binding:"omitempty"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" binding:"omitempty"`
	DeletedAt time.Time `json:"deleted_at" db:"deleted_at" binding:"omitempty"`
}
