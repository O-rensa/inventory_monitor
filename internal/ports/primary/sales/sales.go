package pp_sales

import "time"

type CreateSalesPayload struct {
	SaleDate  time.Time `json:"saleDate" validate:"required"`
	ItemsSold []SaleItemPayload
}

type SaleItemPayload struct {
	ProductId string `json:"productId" validate:"required"`
	ItemCount string `json:"itemCount" validate:"required"`
}
