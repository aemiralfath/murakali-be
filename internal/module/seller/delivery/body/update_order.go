package body

import (
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"
)

type UpdateNoResiOrderSellerRequest struct {
	NoResi string `json:"resi_no"`
}

func (r *UpdateNoResiOrderSellerRequest) ValidateUpdateNoResi() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"resi_no": "",
		},
	}

	r.NoResi = strings.TrimSpace(r.NoResi)
	if r.NoResi == "" {
		unprocessableEntity = true
		entity.Fields["resi_no"] = FieldCannotBeEmptyMessage
	}

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
