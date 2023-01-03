package body





type GetSearchProductQueryRequest struct {
	Search            string
	Sort     string
	SortBy  string
	// array of location
	Category string
	MinPrice float64
	MaxPrice float64
	MinRating float64
	MaxRating float64
}