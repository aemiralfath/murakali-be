package repository

const (
	CheckEmailHistoryQuery  = `SELECT "id", "email" FROM "email_history" WHERE "email" ILIKE $1`
	GetUserByEmailQuery     = `SELECT "id", "email", "is_verify" FROM "user" WHERE "email" ILIKE $1`
	CreateUserQuery         = `INSERT INTO "user" (role_id, email, is_sso, is_verify) VALUES ($1, $2, $3, $4) RETURNING "id", "email"`
	CreateEmailHistoryQuery = `INSERT INTO "email_history" (email) VALUES ($1)`
	VerifyUserQuery         = `UPDATE "user" SET "is_verify" = $1 WHERE "email" = $2`
)
