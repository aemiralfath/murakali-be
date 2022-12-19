package body

import (
	"github.com/google/uuid"
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"
)

type DeleteAddressRequest struct {
	ID string `json:"id"`
}

func (r *DeleteAddressRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"id": "",
		},
	}

	r.ID = strings.TrimSpace(r.ID)
	id, err := uuid.Parse(r.ID)
	if err != nil {
		unprocessableEntity = true
		entity.Fields["id"] = IDNotValidMessage
	}

	r.ID = id.String()
	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
