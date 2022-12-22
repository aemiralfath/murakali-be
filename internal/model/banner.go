package model

import (
	"github.com/google/uuid"
)

type Banner struct {
	ID       uuid.UUID `json:"id" db:"id" binding:"omitempty"`
	Title    string    `json:"title" db:"title" binding:"omitempty"`
	Content  string    `json:"content" db:"content" binding:"omitempty"`
	ImageURL string    `json:"image_url" db:"image_url" binding:"omitempty"`
	PageURL  string    `json:"page_url" db:"page_url" binding:"omitempty"`
	IsActive bool      `json:"is_active" db:"is_active" binding:"omitempty"`
}
