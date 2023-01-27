package body

type SellerPerformance struct {
	ShopID             string                `json:"shop_id" db:"shop_id"`
	ShopName           string                `json:"shop_name" db:"shop_name"`
	ShopCreatedAt      string                `json:"shop_created_at" db:"shop_created_at"`
	ReportUpdatedAt    string                `json:"report_updated_at" db:"report_updated_at"`
	DailySales         []*DailySales         `json:"daily_sales"`
	DailyOrder         []*DailyOrder         `json:"daily_order"`
	MonthlyOrder       *MonthlyOrder         `json:"monthly_order"`
	TotalRating        *TotalRating          `json:"total_rating"`
	MostOrderedProduct []*MostOrderedProduct `json:"most_ordered_product"`
	NumOrderByProvince []*NumOrderByProvince `json:"num_order_by_province"`
	TotalSales         *TotalSales           `json:"total_sales"`
}

type DailySales struct {
	Date       string  `json:"date" db:"date"`
	TotalSales float64 `json:"total_sales" db:"total_sales"`
}

type DailyOrder struct {
	Date         string `json:"date" db:"date"`
	SuccessOrder int    `json:"success_order" db:"success_order"`
	FailedOrder  int    `json:"failed_order" db:"failed_order"`
}

type MonthlyOrder struct {
	Month                     string `json:"month" db:"month"`
	SuccessOrder              int    `json:"success_order" db:"success_order"`
	FailedOrder               int    `json:"failed_order" db:"failed_order"`
	SuccessOrderPercentChange *int   `json:"success_order_percent_change,omitempty" db:"success_order_percent_change"`
	FailedOrderPercentChange  *int   `json:"failed_order_percent_change,omitempty" db:"failed_order_percent_change"`
}

type TotalRating struct {
	Rating1 int `json:"rating_1" db:"rating_1"`
	Rating2 int `json:"rating_2" db:"rating_2"`
	Rating3 int `json:"rating_3" db:"rating_3"`
	Rating4 int `json:"rating_4" db:"rating_4"`
	Rating5 int `json:"rating_5" db:"rating_5"`
}

type MostOrderedProduct struct {
	ID           string `json:"id" db:"id"`
	Title        string `json:"title" db:"title"`
	ViewCount    int    `json:"view_count" db:"view_count"`
	UnitSold     int    `json:"unit_sold" db:"unit_sold"`
	ThumbnailURL string `json:"thumbnail_url" db:"thumbnail_url"`
}

type NumOrderByProvince struct {
	ProvinceID int `json:"province_id" db:"province_id"`
	NumOrders  int `json:"num_orders" db:"num_orders"`
}

type TotalSales struct {
	TotalSales      float64 `json:"total_sales" db:"total_sales"`
	WithdrawableSum float64 `json:"withdrawable_sum" db:"withdrawable_sum"`
	WithdrawnSum    float64 `json:"withdrawn_sum" db:"withdrawn_sum"`
}
