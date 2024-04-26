package as_products

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	c_brand "github.com/o-rensa/iv/internal/core/brands"
	c_inventories "github.com/o-rensa/iv/internal/core/inventory"
	c_productcategory "github.com/o-rensa/iv/internal/core/productcategories"
	c_product "github.com/o-rensa/iv/internal/core/products"
	pp_product "github.com/o-rensa/iv/internal/ports/primary/products"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ProductStore struct {
	db *gorm.DB
}

func NewProductStore(db *sql.DB) (*ProductStore, error) {
	ps := &ProductStore{}

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err == nil {
		ps.db = gdb
	}

	return ps, err
}

// private methods
func (ps *ProductStore) getItemCount(prodID uuid.UUID) (uint, error) {
	var iCount uint = 0
	var inv c_inventories.Inventory
	res := ps.db.First(&inv, "productID = ?", prodID.String())
	if res.Error == nil {
		iCount = inv.ItemCount
	}

	return iCount, res.Error
}

func (ps *ProductStore) getProductCategoryName(pCID uuid.UUID) (string, error) {
	var n string = ""
	var pc c_productcategory.ProductCategory
	res := ps.db.First(&pc, "id = ?", pCID.String())
	if res.Error == nil {
		n = pc.ProductCategoryName
	}

	return n, res.Error
}

func (ps *ProductStore) getBrandName(bID uuid.UUID) (string, error) {
	var n string = ""
	var b c_brand.Brand
	res := ps.db.First(&b, "id = ?", bID.String())
	if res.Error == nil {
		n = b.BrandName
	}

	return n, nil
}

// public methods
func (ps *ProductStore) CreateProduct(input c_product.Product) (pp_product.ProductDto, error) {

	prodID := input.ModelID.ID
	dto := pp_product.ProductDto{}
	inv := c_inventories.Inventory{
		ProductID: prodID,
		ItemCount: 0,
	}

	tx := ps.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return dto, err
	}

	// create product
	if err := tx.Create(&input).Error; err != nil {
		tx.Rollback()
		return dto, err
	}

	// create inventory
	if err := tx.Create(&inv).Error; err != nil {
		tx.Rollback()
		return dto, err
	}

	// get product category name
	pCName, err := ps.getProductCategoryName(input.ProductCategoryID)
	if err != nil {
		tx.Rollback()
		return dto, err
	}

	// get brand name
	bName, err := ps.getBrandName(input.BrandID)
	if err != nil {
		tx.Rollback()
		return dto, err
	}

	// insert values to dto
	dto.ProductID = prodID.String()
	dto.ProductCategoryID = input.ProductCategoryID.String()
	dto.ProductCategoryName = pCName
	dto.BrandID = input.ID.String()
	dto.BrandName = bName
	dto.ProductName = input.ProductName
	dto.ItemCount = "0"

	return dto, tx.Commit().Error
}

func (ps *ProductStore) GetAllProducts() ([]pp_product.ProductDto, error) {
	pDtos := []pp_product.ProductDto{}
	var prods []c_product.Product
	res := ps.db.Find(&prods)
	var retErr error
	if res.Error == nil {
		retErr = res.Error
		for _, v := range prods {
			// get item count
			iCount, err := ps.getItemCount(v.ModelID.ID)
			if err != nil {
				retErr = err
				break
			}

			// get product category name
			pCName, err := ps.getProductCategoryName(v.ProductCategoryID)
			if err != nil {
				retErr = err
				break
			}
			// get brand name
			bName, err := ps.getBrandName(v.BrandID)
			if err != nil {
				retErr = err
				break
			}

			i := pp_product.ProductDto{
				ProductID:           v.ModelID.ID.String(),
				ProductCategoryID:   v.ProductCategoryID.String(),
				ProductCategoryName: pCName,
				BrandID:             v.BrandID.String(),
				BrandName:           bName,
				ProductName:         v.ProductName,
				ItemCount:           fmt.Sprint(iCount),
			}
			pDtos = append(pDtos, i)
		}
	}

	if retErr != nil {
		pDtos = []pp_product.ProductDto{}
	}

	return pDtos, retErr
}

func (ps *ProductStore) GetProductByID(iD uuid.UUID) (pp_product.ProductDto, error) {
	var dto pp_product.ProductDto
	var p c_product.Product
	res := ps.db.First(&p, "id = ?", iD.String())
	if res.Error == nil {
		// get item count
		iCount, err := ps.getItemCount(iD)
		if err != nil {
			return dto, err
		}

		// get product category name
		pName, err := ps.getProductCategoryName(p.ProductCategoryID)
		if err != nil {
			return dto, err
		}

		// get brand name
		bName, err := ps.getBrandName(p.BrandID)
		if err != nil {
			return dto, err
		}

		dto = pp_product.ProductDto{
			ProductID:           p.ModelID.ID.String(),
			ProductCategoryID:   p.ProductCategoryID.String(),
			ProductCategoryName: pName,
			BrandID:             p.BrandID.String(),
			BrandName:           bName,
			ProductName:         p.ProductName,
			ItemCount:           fmt.Sprint(iCount),
		}
	}

	return dto, nil
}

func (ps *ProductStore) UpdateProduct(input c_product.Product) (pp_product.ProductDto, error) {
	var dto pp_product.ProductDto
	var p c_product.Product

	res := ps.db.First(&p, "id = ?", input.ModelID.ID.String())
	if res.Error != nil {
		return dto, res.Error
	}

	updated := ps.db.Model(&p).Updates(map[string]interface{}{
		"productCategoryId": input.ProductCategoryID,
		"brandId":           input.BrandID,
		"productName":       input.ProductName,
	})
	if updated.Error == nil {
		// get item count
		iCount, err := ps.getItemCount(input.ModelID.ID)
		if err != nil {
			return dto, err
		}

		// get product category name
		pCName, err := ps.getProductCategoryName(input.ProductCategoryID)
		if err != nil {
			return dto, err
		}

		// get brand name
		bName, err := ps.getBrandName(input.BrandID)
		if err != nil {
			return dto, err
		}

		dto = pp_product.ProductDto{
			ProductID:           input.ModelID.ID.String(),
			ProductCategoryID:   input.ProductCategoryID.String(),
			ProductCategoryName: pCName,
			BrandID:             input.BrandID.String(),
			BrandName:           bName,
			ProductName:         input.ProductName,
			ItemCount:           fmt.Sprint(iCount),
		}
	}

	return dto, res.Error
}

func (ps *ProductStore) UpdateProductItemCount(prodID uuid.UUID, count uint) error {
	item := ps.db.Model(&c_inventories.Inventory{}).Where("productID = ?", prodID).Update("itemCount", count)
	if item.Error != nil {
		return item.Error
	}

	return nil
}

func (ps *ProductStore) DeleteProduct(iD uuid.UUID) error {
	tx := ps.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	// delete product
	if err := tx.Delete(&c_product.Product{}, iD).Error; err != nil {
		tx.Rollback()
		return err
	}

	// delete inventory
	if err := tx.Delete(&c_inventories.Inventory{}, iD).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
