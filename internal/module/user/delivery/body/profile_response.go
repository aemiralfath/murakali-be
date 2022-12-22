package body

import "time"

type ProfileResponse struct {
	Role        int       `json:"role"`
	UserName    *string   `json:"user_name"`
	Email       string    `json:"email"`
	PhoneNumber *string   `json:"phone_number"`
	FullName    *string   `json:"full_name"`
	Gender      *string   `json:"gender"`
	BirthDate   time.Time `json:"birth_date"`
	PhotoURL    *string   `json:"photo_url"`
	IsVerify    bool      `json:"is_verify"`
}
