package repository

const (
	CheckEmailHistoryQuery  = `SELECT "id", "email" FROM "email_history" WHERE "email" ILIKE $1`
	GetUserByIDQuery        = `SELECT "id", "role_id", "email" FROM "user" WHERE "id" = $1`
	GetUserByEmailQuery     = `SELECT "id", "role_id", "email", "password", "username", "is_verify" FROM "user" WHERE "email" ILIKE $1`
	GetUserByUsernameQuery  = `SELECT "id", "email", "is_verify" FROM "user" WHERE "username" ILIKE $1`
	GetUserByPhoneNoQuery   = `SELECT "id", "email", "is_verify" FROM "user" WHERE "phone_no" ILIKE $1`
	CreateUserQuery         = `INSERT INTO "user" (role_id, email, is_sso, is_verify) VALUES ($1, $2, $3, $4) RETURNING "id", "email"`
	CreateEmailHistoryQuery = `INSERT INTO "email_history" (email) VALUES ($1)`
	VerifyUserQuery         = `UPDATE "user" SET "phone_no" = $1, "fullname" = $2, "username" = $3, "password" = $4, "is_verify" = $5, "updated_at" = $6 WHERE "email" = $7`
	UpdatePasswordQuery     = `UPDATE "user" SET "password" = $1 WHERE "email" = $2`
)
