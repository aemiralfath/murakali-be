package body

import (
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type CreateProductRequest struct {
	Title         string                       `json:"title"`
	Description   string                       `json:"description"`
	Thumbnail     string                       `json:"thumbnail"`
	CategoryID    string                       `json:"category_id"`
	ProductDetail []CreateProductDetailRequest `json:"product_detail"`
}

type CreateProductDetailRequest struct {
	Price           float64  `json:"price"`
	Stock           float64  `json:"stock"`
	Weight          float64  `json:"weight"`
	Size            float64  `json:"size"`
	Hazardous       bool     `json:"hazardous"`
	Codition        string   `json:"condition"`
	BulkPrice       bool     `json:"bulk_price"`
	Photo           []string `json:"photo"`
	Video           []string `json:"video"`
	VariantDetailID []string `json:"variant_detail_id"`
}

func (r *CreateProductRequest) ValidateCreateProduct() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"title":       "",
			"description": "",
			"thumbnail":   "",
			"category_id": "",
		},
	}

	r.Title = strings.TrimSpace(r.Title)
	if r.Title == "" {
		unprocessableEntity = true
		entity.Fields["title"] = FieldCannotBeEmptyMessage
	}

	r.Description = strings.TrimSpace(r.Description)
	if r.Description == "" {
		unprocessableEntity = true
		entity.Fields["description"] = FieldCannotBeEmptyMessage
	}

	r.Thumbnail = strings.TrimSpace(r.Thumbnail)
	if r.Thumbnail == "" {
		unprocessableEntity = true
		entity.Fields["thumbnail"] = FieldCannotBeEmptyMessage
	}

	r.CategoryID = strings.TrimSpace(r.CategoryID)
	if _, err := uuid.Parse(r.CategoryID); err != nil {
		unprocessableEntity = true
		entity.Fields["category_id"] = FieldCannotBeEmptyMessage
	}

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
