package body

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"murakali/config"
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
)

type SLPCallbackRequest struct {
	TxnID        string `json:"txn_id"`
	Amount       string `json:"amount"`
	MerchantCode string `json:"merchant_code"`
	Status       string `json:"status"`
	Message      string `json:"message"`
	Signature    string `json:"signature"`
}

func (r *SLPCallbackRequest) Validate(cfg *config.Config) (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"txn_id":        "",
			"amount":        "",
			"merchant_code": "",
			"status":        "",
			"message":       "",
			"signature":     "",
		},
	}

	signFormat := fmt.Sprintf("%s:%s:%s:%s:%s", r.TxnID, r.Amount, cfg.External.SlpMerchantCode, r.Status, r.Message)
	h := hmac.New(sha256.New, []byte(cfg.External.SlpAPIKey))
	h.Write([]byte(signFormat))
	sign := hex.EncodeToString(h.Sum(nil))

	if sign != r.Signature {
		unprocessableEntity = true
		entity.Fields["txn_id"] = r.TxnID
		entity.Fields["amount"] = r.Amount
		entity.Fields["merchant_code"] = r.MerchantCode
		entity.Fields["status"] = r.Status
		entity.Fields["message"] = r.Message
		entity.Fields["signature"] = InvalidSignatureMessage
	}

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
