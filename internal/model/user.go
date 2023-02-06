package model

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID    `json:"id" db:"id" binding:"omitempty"`
	RoleID    int          `json:"role_id" db:"role_id" binding:"omitempty"`
	Email     string       `json:"email" db:"email" binding:"omitempty"`
	Password  *string      `json:"password" db:"password" binding:"omitempty"`
	Username  *string      `json:"username" db:"username" binding:"omitempty"`
	PhoneNo   *string      `json:"phone_no" db:"phone_no" binding:"omitempty"`
	FullName  *string      `json:"fullname" db:"fullname" binding:"omitempty"`
	Gender    *string      `json:"gender" db:"gender" binding:"omitempty"`
	PhotoURL  *string      `json:"photo_url" db:"photo_url" binding:"omitempty"`
	BirthDate sql.NullTime `json:"birth_date" db:"birth_date" binding:"omitempty"`
	IsSSO     bool         `json:"is_sso" db:"is_sso" binding:"omitempty"`
	IsVerify  bool         `json:"is_verify" db:"is_verify" binding:"omitempty"`
	CreatedAt time.Time    `json:"created_at" db:"created_at" binding:"omitempty"`
	UpdatedAt sql.NullTime `json:"updated_at" db:"updated_at" binding:"omitempty"`
	DeletedAt sql.NullTime `json:"deleted_at" db:"deleted_at" binding:"omitempty"`
}
