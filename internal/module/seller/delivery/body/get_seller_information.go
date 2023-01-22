package body

import (
	"time"

	"github.com/google/uuid"
)

type SellerInformationResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	TotalProduct int       `json:"total_product"`
	TotalRating  float64   `json:"total_rating"`
	RatingAVG    float64   `json:"rating_avg"`
	PhotoURL     string    `json:"photo_url"`
	CreatedAt    time.Time `json:"created_at"`
}
