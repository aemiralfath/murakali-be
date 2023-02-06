package body

import (
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type CategoryRequest struct {
	NameLevelOne   string `form:"name_lvl_one"`
	NameLevelTwo   string `form:"name_lvl_two"`
	NameLevelThree string `form:"name_lvl_three"`
}

type CategoryResponse struct {
	ID       uuid.UUID `json:"id"`
	ParentID uuid.UUID `json:"parent_id"`
	Name     string    `json:"name"`
	PhotoURL string    `json:"photo_url"`

	ChildCategory []*CategoryResponse `json:"child_category"`
}

func (r *CategoryRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"name": "",
		},
	}

	r.NameLevelOne = strings.TrimSpace(r.NameLevelOne)
	if r.NameLevelOne == "" {
		unprocessableEntity = true
		entity.Fields["name_lvl_one"] = FieldCannotBeEmptyMessage
	}

	r.NameLevelTwo = strings.TrimSpace(r.NameLevelTwo)
	if r.NameLevelTwo == "" {
		unprocessableEntity = true
		entity.Fields["name_lvl_two"] = FieldCannotBeEmptyMessage
	}

	r.NameLevelThree = strings.TrimSpace(r.NameLevelThree)
	if r.NameLevelThree == "" {
		unprocessableEntity = true
		entity.Fields["name_lvl_three"] = FieldCannotBeEmptyMessage
	}
	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
