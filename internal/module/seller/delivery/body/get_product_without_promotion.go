package body

type GetProductWithoutPromotion struct {
	ProductID               string       `json:"product_id"`
	ProductName             string       `json:"product_name"`
	Price                float64      `json:"price"`
	CategoryName string `json:"category_name"`
	ProductThumbnailURL     string       `json:"product_thumbnail_url"`
	UnitSold int `json:"unit_sold"`
	Rating float64 `json:"rating"`
	
}