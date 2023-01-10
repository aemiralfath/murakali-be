package body

import (
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type UpdateProductRequest struct {
	ProductInfo   UpdateProductInfo            `json:"products_info"`
	ProductDetail []UpdateProductDetailRequest `json:"products_detail"`
}

type UpdateProductInfo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
	CategoryID  string `json:"category_id"`
}

type UpdateProductInfoForQuery struct {
	Title       string
	Description string
	Thumbnail   string
	CategoryID  string
	MinPrice    float64
	MaxPrice    float64
}

type UpdateProductDetailRequest struct {
	ProductDetailID string          `json:"product_detail_id"`
	Price           float64         `json:"price"`
	Stock           float64         `json:"stock"`
	Weight          float64         `json:"weight"`
	Size            float64         `json:"size"`
	Hazardous       bool            `json:"hazardous"`
	Codition        string          `json:"condition"`
	BulkPrice       bool            `json:"bulk_price"`
	Photo           []string        `json:"photo"`
	VariantDetailID []UpdateVariant `json:"variant_info"`
}

type UpdateVariant struct {
	VariantID       string `json:"variant_id"`
	VariantDetailID string `json:"variant_detail_id"`
}

func (r *UpdateProductRequest) ValidateUpdateProduct() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"products_info.title":       "",
			"products_info.description": "",
			"products_info.thumbnail":   "",
			"products_info.category_id": "",
		},
	}

	r.ProductInfo.Title = strings.TrimSpace(r.ProductInfo.Title)
	if r.ProductInfo.Title == "" {
		unprocessableEntity = true
		entity.Fields["products_info.title"] = FieldCannotBeEmptyMessage
	}

	r.ProductInfo.Description = strings.TrimSpace(r.ProductInfo.Description)
	if r.ProductInfo.Description == "" {
		unprocessableEntity = true
		entity.Fields["products_info.description"] = FieldCannotBeEmptyMessage
	}

	r.ProductInfo.Thumbnail = strings.TrimSpace(r.ProductInfo.Thumbnail)
	if r.ProductInfo.Thumbnail == "" {
		unprocessableEntity = true
		entity.Fields["products_info.thumbnail"] = FieldCannotBeEmptyMessage
	}

	r.ProductInfo.CategoryID = strings.TrimSpace(r.ProductInfo.CategoryID)
	if _, err := uuid.Parse(r.ProductInfo.CategoryID); err != nil {
		unprocessableEntity = true
		entity.Fields["products_info.category_id"] = FieldCannotBeEmptyMessage
	}

	totalData := len(r.ProductDetail)
	if totalData == 0 {
		unprocessableEntity = true
		entity.Fields["products_detail"] = FieldCannotBeEmptyMessage
	}
	for i := 0; i < totalData; i++ {
		r.ProductDetail[i].ProductDetailID = strings.TrimSpace(r.ProductDetail[i].ProductDetailID)
		if r.ProductDetail[i].ProductDetailID == "" {
			unprocessableEntity = true
			entity.Fields["product_detail_id"] = FieldCannotBeEmptyMessage
		}
		if r.ProductDetail[i].Price == 0 {
			unprocessableEntity = true
			entity.Fields["price"] = FieldCannotBeEmptyMessage
		}
		if r.ProductDetail[i].Stock == 0 {
			unprocessableEntity = true
			entity.Fields["stock"] = FieldCannotBeEmptyMessage
		}
		if r.ProductDetail[i].Weight == 0 {
			unprocessableEntity = true
			entity.Fields["weight"] = FieldCannotBeEmptyMessage
		}
		if r.ProductDetail[i].Size == 0 {
			unprocessableEntity = true
			entity.Fields["size"] = FieldCannotBeEmptyMessage
		}
		r.ProductDetail[i].Codition = strings.TrimSpace(r.ProductDetail[i].Codition)
		if r.ProductDetail[i].Codition == "" {
			unprocessableEntity = true
			entity.Fields["condition"] = FieldCannotBeEmptyMessage
		}
		if len(r.ProductDetail[i].Photo) == 0 {
			unprocessableEntity = true
			entity.Fields["photo"] = FieldCannotBeEmptyMessage
		}

	}

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
