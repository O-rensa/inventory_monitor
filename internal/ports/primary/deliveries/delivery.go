package pp_delivery

import "time"

type CreateDeliveryPayload struct {
	DeliveryDate   time.Time `json:"deliveryDate" validate:"required"`
	ItemsDelivered []DeliveryItemPayload
}

type DeliveryItemPayload struct {
	ProductId string `json:"productId" validate:"required"`
	ItemCount string `json:"itemCount" validate:"required"`
}
