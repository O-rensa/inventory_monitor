package ps_sales

import (
	"github.com/google/uuid"
	c_sales "github.com/o-rensa/iv/internal/core/sales"
)

type SalesStore interface {
	CreateNewSalesTransaction(c_sales.SaleDetails) error
	DeleteSalesTransaction(uuid.UUID) error
}
