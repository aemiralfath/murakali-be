package body

const (
	FieldCannotBeEmptyMessage = "Field cannot be empty."
	ProductNotFound           = "Product not found"
	UpdateProductFailed       = "Update product failed"
	ImageIsEmpty              = "image cannot be empty"
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
	Province  []string
}

type GetImageResponse struct {
	ProductDetailId *string `json:"product_detail_id"`
	Url             string  `json:"url"`
}

type GetAllProductImageResponse struct {
	Image []*GetImageResponse `json:"image"`
}
