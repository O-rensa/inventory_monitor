package ps_productCategory

import (
	c_productcategory "github.com/o-rensa/iv/internal/core/productcategories"
	pp_productCategory "github.com/o-rensa/iv/internal/ports/primary/productcategories"
)

type ProductCategoryStore interface {
	CreateProductCategory(c_productcategory.ProductCategory) (pp_productCategory.ProductCategoryDto, error)
}
