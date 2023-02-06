package body

import (
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	ReviewAlreadyExist = "Product Review Already Exist"
	ReviewNotExist     = "Review Not Exist"
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
	UserID      string `json:"user_id"`
}

func (p *GetReviewQueryRequest) GetValidate() string {
	query := ""

	if p.UserID != "" {
		query += ` AND r.user_id = '` + p.UserID + `'`
	}

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

type ReviewProductRequest struct {
	ProductID string  `json:"product_id"`
	Comment   *string `json:"comment,omitempty"`
	Rating    int     `json:"rating"`
	PhotoURL  *string `json:"photo_url,omitempty"`
}

func (r *ReviewProductRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"product_id": "",
			"rating":     "",
		},
	}

	r.ProductID = strings.TrimSpace(r.ProductID)
	if r.ProductID == "" {
		unprocessableEntity = true
		entity.Fields["product_id"] = FieldCannotBeEmptyMessage
	}

	if r.Rating == 0 {
		unprocessableEntity = true
		entity.Fields["rating"] = FieldCannotBeEmptyMessage
	}

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}

type DeleteReviewProductRequest struct {
	ReviewID string `json:"review_id"`
}

func (r *DeleteReviewProductRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"review_id": "",
		},
	}

	r.ReviewID = strings.TrimSpace(r.ReviewID)
	if r.ReviewID == "" {
		unprocessableEntity = true
		entity.Fields["product_id"] = FieldCannotBeEmptyMessage
	}

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
