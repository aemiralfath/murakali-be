package body

const (
	FieldCannotBeEmptyMessage = "Field cannot be empty."
)

type UnprocessableEntity struct {
	Fields map[string]string `json:"fields"`
}

type GetProductQueryRequest struct {
	Search    string
	Sort      string
	SortBy    string
	Shop      string
	Category  string
	MinPrice  float64
	MaxPrice  float64
	MinRating float64
	MaxRating float64
	Province  string
}
