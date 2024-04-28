package ps_product

import (
	"github.com/google/uuid"
	c_product "github.com/o-rensa/iv/internal/core/products"
	pp_product "github.com/o-rensa/iv/internal/ports/primary/products"
)

type ProductStore interface {
	CreateProduct(c_product.Product) (pp_product.ProductDto, error)
	GetAllProducts() ([]pp_product.ProductDto, error)
	GetProductByID(uuid.UUID) (pp_product.ProductDto, error)
	UpdateProduct(c_product.Product) (pp_product.ProductDto, error)
	UpdateProductItemCount(uuid.UUID, uint64) error
	DeleteProduct(uuid.UUID) error
}
