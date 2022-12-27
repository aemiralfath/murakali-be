package body

type GetAddressDefaultRequest struct {
	IsDefault   string `form:"is_default"`
	ShopDefault string `form:"shop_default"`
}
