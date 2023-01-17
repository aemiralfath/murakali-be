package body

import (
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"
	"time"
)

type UpdateNoResiOrderSellerRequest struct {
	NoResi               string `json:"resi_no"`
	EstimateArriveAt     string `json:"estimate_arrive_at"`
	EstimateArriveAtTime time.Time
}

func (r *UpdateNoResiOrderSellerRequest) ValidateUpdateNoResi() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"resi_no":            "",
			"estimate_arrive_at": "",
		},
	}

	r.NoResi = strings.TrimSpace(r.NoResi)
	if r.NoResi == "" {
		unprocessableEntity = true
		entity.Fields["resi_no"] = FieldCannotBeEmptyMessage
	}

	r.EstimateArriveAt = strings.TrimSpace(r.EstimateArriveAt)
	estimateTime, err := time.Parse("02-01-2006 15:04:05", r.EstimateArriveAt)
	if err != nil {
		unprocessableEntity = true
		entity.Fields["estimate_arrive_at"] = InvalidDateFormatMessage
	}

	r.EstimateArriveAtTime = estimateTime

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
