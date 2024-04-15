package ps_productCategory

import (
	"github.com/google/uuid"
	c_productcategory "github.com/o-rensa/iv/internal/core/productcategories"
	pp_productCategory "github.com/o-rensa/iv/internal/ports/primary/productcategories"
)

type ProductCategoryStore interface {
	CreateProductCategory(c_productcategory.ProductCategory) (pp_productCategory.ProductCategoryDto, error)
	GetAllProductCategory() ([]pp_productCategory.ProductCategoryDto, error)
	GetProductCategoryByID(uuid.UUID) (pp_productCategory.ProductCategoryDto, error)
	UpdateProductCategory(c_productcategory.ProductCategory) (pp_productCategory.ProductCategoryDto, error)
	DeleteProductCategory(uuid.UUID) error
}
