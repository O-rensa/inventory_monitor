package tests

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

func TestCreateProductCategoryHandler(t *testing.T) {
	apiUrl := "/create"

	// cases
	noProductCategoryNameCase := []byte(`{"productCategoryName":}`)
	validCase := []byte(`{"productCategoryName": "SuperProductCategory"}`)

	t.Run("should fail because productCategoryName is missing", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(noProductCategoryNameCase))
		if err != nil {
			t.Fatal(err)
		}
		defer req.Body.Close()

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc(apiUrl, handler.CreateProductCategoryHandler).Methods(http.MethodPost)

		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Error(rr.Body)
			t.Errorf("expected status cure %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should succedd cause payload is valid", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(validCase))
		if err != nil {
			t.Fatal(err)
		}
		defer req.Body.Close()

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc(apiUrl, handler.CreateProductCategoryHandler).Methods(http.MethodPost)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Error(rr.Body)
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})
}

func TestGetAllProductCategoryHandler(t *testing.T) {
	apiUrl := "/getall"

	t.Run("should succeed", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, apiUrl, strings.NewReader(""))
		if err != nil {
			t.Fatal(err)
		}
		defer req.Body.Close()

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc(apiUrl, handler.GetAllProductCategoryHandler).Methods(http.MethodGet)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Error(rr.Body)
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})
}

func TestGetProductCategoryByIDHandler(t *testing.T) {
	wrongApiUrl := "/get/" + wrongUUID
	corrApiUrl := "/get/" + mockUUID.String()
	apiUrl := "/get/{iD}"

	t.Run("should fail because the uuid is not correct", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, wrongApiUrl, strings.NewReader(""))
		if err != nil {
			t.Fatal(err)
		}
		defer req.Body.Close()

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc(apiUrl, handler.GetProductCategoryByIDHandler).Methods(http.MethodGet)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Error(rr.Body)
			t.Errorf("expected status code %d, got %d", http.StatusNotFound, rr.Code)
		}
	})

	t.Run("should succedd because the uuid is correct", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, corrApiUrl, strings.NewReader(""))
		if err != nil {
			t.Fatal(err)
		}
		defer req.Body.Close()

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc(apiUrl, handler.GetProductCategoryByIDHandler).Methods(http.MethodGet)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Error(rr.Body)
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})
}

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
	if i != mockUUID {
		return dto, errors.New("id did not match to mockUUID")
	}
	dto.ProductCategoryID = mockUUID.String()
	dto.ProductCategoryName = "Product Category Uno"
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
