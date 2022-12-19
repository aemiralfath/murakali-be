package repository

const (
	GetSealabsPayByIdQuery        = `SELECT * from sealabs_pay where user_id = $1`
)
