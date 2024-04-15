package as_productCategory

import (
	"database/sql"

	"github.com/google/uuid"
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

func (pcs *ProductCategoryStore) GetAllProductCategory() ([]pp_productCategory.ProductCategoryDto, error) {
	var dto []pp_productCategory.ProductCategoryDto
	var pc []c_productcategory.ProductCategory
	res := pcs.db.Find(&pc)
	if res.Error == nil {
		for _, v := range pc {
			i := pp_productCategory.ProductCategoryDto{
				ProductCategoryID:   v.ModelID.ID.String(),
				ProductCategoryName: v.ProductCategoryName,
			}
			dto = append(dto, i)
		}
	}
	return dto, nil
}

func (pcs *ProductCategoryStore) GetProductCategoryByID(iD uuid.UUID) (pp_productCategory.ProductCategoryDto, error) {
	var dto pp_productCategory.ProductCategoryDto
	var pc c_productcategory.ProductCategory
	res := pcs.db.First(&pc, "id = ?", iD.String())
	if res.Error == nil {
		dto = pp_productCategory.ProductCategoryDto{
			ProductCategoryID:   pc.ModelID.ID.String(),
			ProductCategoryName: pc.ProductCategoryName,
		}
	}

	return dto, res.Error
}

func (pcs *ProductCategoryStore) UpdateProductCategory(input c_productcategory.ProductCategory) (pp_productCategory.ProductCategoryDto, error) {
	var dto pp_productCategory.ProductCategoryDto
	var pc c_productcategory.ProductCategory

	gdb := pcs.db.First(&pc, "id = ?", input.ModelID.ID.String())
	if gdb.Error != nil {
		return dto, gdb.Error
	}

	res := pcs.db.Model(&pc).Updates(map[string]interface{}{
		"productCategoryName": input.ProductCategoryName,
	})
	if res.Error == nil {
		dto = pp_productCategory.ProductCategoryDto{
			ProductCategoryID:   input.ID.String(),
			ProductCategoryName: input.ProductCategoryName,
		}
	}
	return dto, res.Error
}

func (pcs *ProductCategoryStore) DeleteProductCategory(iD uuid.UUID) error {
	return (pcs.db.Delete(&c_productcategory.ProductCategory{}, iD)).Error
}
