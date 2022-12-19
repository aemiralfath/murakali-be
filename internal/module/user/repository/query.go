package repository

const (
	GetSealabsPayByIdQuery         = `SELECT * from sealabs_pay where user_id = $1`
	CreateSealabsPayQuery          = `INSERT INTO "sealabs_pay" (card_number, user_id, name, is_default,active_date) VALUES ($1, $2, $3, $4, $5)`
	CheckDefaultSealabsPayQuery    = `SELECT card_number from "sealabs_pay" where user_id = $1 and is_default is true`
	SetDefaultSealabsPayTransQuery = `UPDATE "sealabs_pay" set is_default = FALSE where card_number = $1`
	PatchSealabsPayQuery           = `UPDATE "sealabs_pay" set is_default = TRUE where card_number = $1`
	SetDefaultSealabsPayQuery      = `UPDATE "sealabs_pay" set is_default = FALSE where card_number <> $1 and user_id = $2`
)
