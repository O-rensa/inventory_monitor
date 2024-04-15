package tests

import (
	"errors"

	"github.com/google/uuid"
	ap_productCategory "github.com/o-rensa/iv/internal/adapters/primary/productCategories"
	c_admin "github.com/o-rensa/iv/internal/core/admin"
	c_productcategory "github.com/o-rensa/iv/internal/core/productcategories"
	pp_productCategory "github.com/o-rensa/iv/internal/ports/primary/productcategories"
)

var productCategoryStore = &mockProductCategoryStore{}
var adminStore = &mockAdminStore{}
var handler = ap_productCategory.NewProductCategoryHandler(productCategoryStore, adminStore)
var mockUUID = uuid.MustParse("56bd6b24-83e5-4da7-ab26-e45cbe0aa8b6")
var wrongUUID = "7e742f61-58ed-4072-9980-b4170176ca44"

// mock product category store
type mockProductCategoryStore struct{}

func (m *mockProductCategoryStore) CreateProductCategory(i c_productcategory.ProductCategory) (pp_productCategory.ProductCategoryDto, error) {
	dto := pp_productCategory.ProductCategoryDto{}
	dto.ProductCategoryID = i.ModelID.ID.String()
	dto.ProductCategoryName = i.ProductCategoryName
	return dto, nil
}

func (m *mockProductCategoryStore) GetAllProductCategory() ([]pp_productCategory.ProductCategoryDto, error) {
	dtos := []pp_productCategory.ProductCategoryDto{}
	return dtos, nil
}

func (m *mockProductCategoryStore) GetProductCategoryByID(i uuid.UUID) (pp_productCategory.ProductCategoryDto, error) {
	dto := pp_productCategory.ProductCategoryDto{}
	return dto, nil
}

func (m *mockProductCategoryStore) UpdateProductCategory(i c_productcategory.ProductCategory) (pp_productCategory.ProductCategoryDto, error) {
	dto := pp_productCategory.ProductCategoryDto{}
	return dto, nil
}

func (m *mockProductCategoryStore) DeleteProductCategory(i uuid.UUID) error {
	return nil
}

// mock admin store
type mockAdminStore struct{}

func (m *mockAdminStore) CreateAdmin(a c_admin.Admin) error {
	return nil
}

func (m *mockAdminStore) GetAdminByUsername(u string) (c_admin.Admin, error) {
	av := c_admin.Admin{}
	return av, errors.New("error")
}

func (m *mockAdminStore) GetAdminByID(iD uuid.UUID) (c_admin.Admin, error) {
	av := c_admin.Admin{}
	return av, nil
}

func (m *mockAdminStore) UpdateAdminPassword(iD uuid.UUID, np string) error {
	return nil
}
