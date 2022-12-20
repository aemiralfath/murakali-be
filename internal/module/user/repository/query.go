package repository

const (
	GetSealabsPayByIdQuery         = `SELECT * from sealabs_pay where user_id = $1 and deleted_at is null`
	CreateSealabsPayQuery          = `INSERT INTO "sealabs_pay" (card_number, user_id, name, is_default,active_date) VALUES ($1, $2, $3, $4, $5)`
	CheckDefaultSealabsPayQuery    = `SELECT card_number from "sealabs_pay" where user_id = $1 and is_default is true and deleted_at is null`
	SetDefaultSealabsPayTransQuery = `UPDATE "sealabs_pay" set is_default = FALSE,updated_at = now() where card_number = $1`
	PatchSealabsPayQuery           = `UPDATE "sealabs_pay" set is_default = TRUE,updated_at = now() where card_number = $1`
	SetDefaultSealabsPayQuery      = `UPDATE "sealabs_pay" set is_default = FALSE where card_number <> $1 and user_id = $2`
	DeleteSealabsPayQuery          = `UPDATE "sealabs_pay" set deleted_at = now() where card_number = $1 and is_default = FALSE`
	GetUserByIDQuery               = `SELECT "id", "role_id", "email", "username", "phone_no", "fullname", "gender", "birth_date", "is_verify" FROM "user" WHERE "id" = $1`
	CheckEmailHistoryQuery         = `SELECT "id", "email" FROM "email_history" WHERE "email" ILIKE $1`
	GetUserByUsernameQuery         = `SELECT "id", "email", "username", "is_verify" FROM "user" WHERE "username" ILIKE $1`
	GetUserByPhoneNoQuery          = `SELECT "id", "email", "phone_no", "is_verify" FROM "user" WHERE "phone_no" ILIKE $1`
	UpdateUserFieldQuery           = `UPDATE "user" SET "username" = $1, "fullname" = $2, "phone_no" = $3, "birth_date" = $4, "gender" = $5, "updated_at" = $6 WHERE "email" = $7`
	UpdateUserEmailQuery           = `UPDATE "user" SET "email" = $1, "updated_at" = $2 WHERE "id" = $3`
	CreateEmailHistoryQuery        = `INSERT INTO "email_history" (email) VALUES ($1)`
)
