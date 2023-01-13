package model

type OrderDetail struct {
	ProductDetailID  string            `json:"product_detail_id"`
	ProductID        string            `json:"product_id"`
	ProductTitle     string            `json:"product_title"`
	ProductDetailURL *string           `json:"product_detail_url"`
	OrderQuantity    int               `json:"order_quantity"`
	ItemPrice        *float64          `json:"order_item_price"`
	TotalPrice       *float64          `json:"order_total_price"`
	Variant          map[string]string `json:"variant"`
}
