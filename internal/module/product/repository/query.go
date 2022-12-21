package repository

const (
	GetCategoriesQuery           = `SELECT "id", "parent_id", "name", "photo_url" FROM "category" WHERE "parent_id" IS NULL AND "deleted_at" IS NULL`
	GetCategoriesByNameQuery     = `SELECT "id", "parent_id", "name", "photo_url" FROM "category" WHERE "name" = $1 AND "deleted_at" IS NULL`
	GetCategoriesByParentIdQuery = `SELECT "id", "parent_id", "name", "photo_url" FROM "category" WHERE "parent_id" = $1 AND "deleted_at" IS NULL`
)
