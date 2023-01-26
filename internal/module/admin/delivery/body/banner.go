package body

import (
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"
)

type BannerRequest struct {
	Title    string `json:"title" `
	Content  string `json:"content"`
	ImageURL string `json:"image_url"`
	PageURL  string `json:"page_url"`
	IsActive bool   `json:"is_active"`
}

type BannerIDRequest struct {
	ID       string `json:"id"`
	IsActive bool   `json:"is_active"`
}

type BannerResponse struct {
	ID       string `json:"id"`
	Title    string `json:"title" `
	Content  string `json:"content"`
	ImageURL string `json:"image_url"`
	PageURL  string `json:"page_url"`
	IsActive bool   `json:"is_active"`
}

func (r *BannerRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"title":     "",
			"content":   "",
			"image_url": "",
			"page_url":  "",
			"isActive":  "",
		},
	}

	r.Title = strings.TrimSpace(r.Title)
	if r.Title == "" {
		unprocessableEntity = true
		entity.Fields["title"] = FieldCannotBeEmptyMessage
	}

	r.Content = strings.TrimSpace(r.Content)
	if r.Content == "" {
		unprocessableEntity = true
		entity.Fields["content"] = FieldCannotBeEmptyMessage
	}

	r.ImageURL = strings.TrimSpace(r.ImageURL)
	if r.ImageURL == "" {
		unprocessableEntity = true
		entity.Fields["image_url"] = FieldCannotBeEmptyMessage
	}

	r.PageURL = strings.TrimSpace(r.PageURL)
	if r.PageURL == "" {
		unprocessableEntity = true
		entity.Fields["page_url"] = FieldCannotBeEmptyMessage
	}

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
