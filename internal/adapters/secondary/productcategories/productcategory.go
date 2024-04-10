package as_productCategory

import (
	"database/sql"

	c_productcategory "github.com/o-rensa/iv/internal/core/productcategories"
	pp_productCategory "github.com/o-rensa/iv/internal/ports/primary/productcategories"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ProductCategoryStore struct {
	db *gorm.DB
}

func NewProductCategoryStore(db *sql.DB) (*ProductCategoryStore, error) {
	pcs := &ProductCategoryStore{}

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err == nil {
		pcs.db = gdb
	}

	return pcs, err
}

func (pcs *ProductCategoryStore) CreateProductCategory(input c_productcategory.ProductCategory) (pp_productCategory.ProductCategoryDto, error) {
	var dto pp_productCategory.ProductCategoryDto
	res := pcs.db.Create(input)
	if res.Error == nil {
		dto = pp_productCategory.ProductCategoryDto{
			ProductCategoryID:   input.ID.String(),
			ProductCategoryName: input.ProductCategoryName,
		}
	}

	return dto, res.Error
}
