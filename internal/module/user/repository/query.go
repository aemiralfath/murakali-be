package repository

const (
	GetUserByIDQuery       = `SELECT "id", "role_id", "email", "username", "phone_no", "fullname", "gender", "birth_date", "is_verify" FROM "user" WHERE "id" = $1`
	GetUserByUsernameQuery = `SELECT "id", "email", "username", "is_verify" FROM "user" WHERE "username" ILIKE $1`
	GetUserByPhoneNoQuery  = `SELECT "id", "email", "phone_no", "is_verify" FROM "user" WHERE "phone_no" ILIKE $1`
	UpdateUserFieldQuery   = `UPDATE "user" SET "username" = $1, "fullname" = $2, "phone_no" = $3, "birth_date" = $4, "gender" = $5, "updated_at" = $6 WHERE "email" = $7`
)
