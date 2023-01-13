package body

import (
	"time"

	"github.com/google/uuid"
)

type ReviewProduct struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	ProductID uuid.UUID `json:"product_id"`
	Comment   *string   `json:"comment"`
	Rating    int       `json:"rating"`
	ImageURL  *string   `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	PhotoURL  *string   `json:"photo_url"`
	Username  string    `json:"username"`
}

type RatingProduct struct {
	Rating int `json:"rating"`
	Count  int `json:"count"`
}

type AllRatingProduct struct {
	TotalRating   float64          `json:"total_rating"`
	AvgRating     float64          `json:"avg_rating"`
	RatingProduct []*RatingProduct `json:"rating_product"`
}

type GetReviewQueryRequest struct {
	Rating      string `json:"rating"`
	ShowComment bool   `json:"show_comment"`
	ShowImage   bool   `json:"show_image"`
}

func (p *GetReviewQueryRequest) GetValidate() string {
	query := ""

	if p.Rating != "" && p.Rating != "0" {
		query += " AND r.rating = " + p.Rating
	}

	if !p.ShowComment && !p.ShowImage {
		query += " AND (r.comment IS NULL OR r.comment = '') AND (r.image_url IS NULL OR r.image_url = '')"
	} else if p.ShowComment && !p.ShowImage {
		query += " AND r.comment IS NOT NULL AND (r.image_url IS NULL OR r.image_url = '')"
	} else if !p.ShowComment && p.ShowImage {
		query += " AND r.image_url IS NOT NULL AND (r.comment IS NULL OR r.comment = '')"
	}

	return query
}
