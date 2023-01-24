package body

import (
	"database/sql"
	"mime/multipart"
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"
)

type ImageRequest struct {
	Img multipart.File `form:"file"`
}

type CategoryRequest struct {
	ID            string `json:"id"`
	ParentID      string `json:"parent_id" `
	ParentIDValue sql.NullString
	Name          string `json:"name" `
	PhotoURL      string `json:"photo_url"`
	Level         string `json:"level" `
}

type CategoryResponse struct {
	CategoryID string  `json:"id" db:"id"`
	ParentID   *string `json:"parent_id" db:"parent_id"`
	Name       string  `json:"name" db:"name"`
	PhotoURL   *string `json:"photo_url" db:"photo_url"`
	Level      string  `json:"level" db:"level"`
}

func (r *CategoryRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"parent_id": "",
			"name":      "",
			"photo_url": "",
			"level":     "",
		},
	}

	r.Name = strings.TrimSpace(r.Name)
	if r.Name == "" {
		unprocessableEntity = true
		entity.Fields["name"] = FieldCannotBeEmptyMessage
	}

	r.ParentID = strings.TrimSpace(r.ParentID)
	if r.ParentID == "" {
		r.ParentIDValue.Valid = false
		r.ParentIDValue.String = ""
	} else {
		r.ParentIDValue.String = r.ParentID
		r.ParentIDValue.Valid = true
	}

	r.PhotoURL = strings.TrimSpace(r.PhotoURL)
	if r.PhotoURL == "" {
		unprocessableEntity = true
		entity.Fields["photo_url"] = FieldCannotBeEmptyMessage
	}

	r.Level = strings.TrimSpace(r.Level)
	if r.Level == "" {
		unprocessableEntity = true
		entity.Fields["level"] = FieldCannotBeEmptyMessage
	}
	if r.Level != "1" && r.Level != "2" && r.Level != "3" {
		unprocessableEntity = true
		entity.Fields["level"] = CategoryLevelInvalid
	}

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
