package ps_delivery

import (
	"github.com/google/uuid"
	c_deliveries "github.com/o-rensa/iv/internal/core/deliveries"
)

type DeliveryStore interface {
	CreateDeliveryTransaction(c_deliveries.DeliveryDetails) error
	DeleteDeliveryTransaction(uuid.UUID) error
}
