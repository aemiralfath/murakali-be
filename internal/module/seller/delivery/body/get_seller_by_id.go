package body

import (
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type SellerByIDRequest struct {
	SellerID string `json:"seller_id"`
}

type SellerResponse struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"-"`
	Name         string    `json:"name"`
	TotalProduct int       `json:"total_product"`
	TotalRating  float64   `json:"total_rating"`
	RatingAVG    float64   `json:"rating_avg"`
	PhotoURL     string    `json:"photo_url"`
	CreatedAt    time.Time `json:"created_at"`
}

func (r *SellerByIDRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"seller_id": "",
		},
	}

	r.SellerID = strings.TrimSpace(r.SellerID)
	if r.SellerID == "" {
		unprocessableEntity = true
		entity.Fields["seller_id"] = FieldCannotBeEmptyMessage
	}

	_, err := uuid.Parse(r.SellerID)
	if err != nil {
		unprocessableEntity = true
		entity.Fields["seller_id"] = IDNotValidMessage
	}

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
