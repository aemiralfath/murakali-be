package body

import (
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"
)

const (
	FieldCannotBeEmptyMessage = "Field cannot be empty."
	ProductNotFound           = "Product not found"
	UpdateProductFailed       = "Update product failed"
	ImageIsEmpty              = "image cannot be empty"
)

type UnprocessableEntity struct {
	Fields map[string]string `json:"fields"`
}

type GetProductQueryRequest struct {
	Search       string
	Sort         string
	SortBy       string
	Shop         string
	Category     string
	MinPrice     float64
	MaxPrice     float64
	MinRating    float64
	MaxRating    float64
	ListedStatus int
	Province     []string
}

type GetProductRequest struct {
	ProductID string `json:"product_id" binding:"required"`
}

func (r *GetProductRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"product_id": "",
		},
	}

	r.ProductID = strings.TrimSpace(r.ProductID)
	if r.ProductID == "" {
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

type GetImageResponse struct {
	ProductDetailID *string `json:"product_detail_id"`
	URL             string  `json:"url"`
}

type GetAllProductImageResponse struct {
	Image []*GetImageResponse `json:"image"`
}
