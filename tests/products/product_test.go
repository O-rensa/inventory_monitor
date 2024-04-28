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
	ap_product "github.com/o-rensa/iv/internal/adapters/primary/product"
	c_admin "github.com/o-rensa/iv/internal/core/admin"
	c_product "github.com/o-rensa/iv/internal/core/products"
	pp_product "github.com/o-rensa/iv/internal/ports/primary/products"
)

var adminStore = &mockAdminStore{}
var productStore = &mockProductStore{}
var handler = ap_product.NewProductHandler(productStore, adminStore)
var mockUUID = uuid.MustParse("56bd6b24-83e5-4da7-ab26-e45cbe0aa8b6")
var wrongUUID = "7e742f61-58ed-4072-9980-b4170176ca44"

// test create product handler
func TestCreateProductHandler(t *testing.T) {
	apiUrl := "/create"

	// cases
	invalidCase := []byte(`{"productCategoryID":, "brandID":, "ProductName":"test product name"}`)
	validCase := []byte(`{"productCategoryID": "4250ad37-e502-4459-bbb0-5ab93094f2bd", "brandID": "012085a4-76ac-4bd6-9f5e-3830e693e902", "productName": "test product name"}`)

	t.Run("should fail because invalid", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(invalidCase))
		if err != nil {
			t.Fatal(err)
		}
		defer req.Body.Close()

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc(apiUrl, handler.CreateProductHandler).Methods(http.MethodPost)

		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Error(rr.Body)
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("shoudl succeed cause payload is valid", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(validCase))
		if err != nil {
			t.Fatal(err)
		}
		defer req.Body.Close()

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc(apiUrl, handler.CreateProductHandler).Methods(http.MethodPost)

		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusCreated {
			t.Error(rr.Body)
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})
}

func TestGetProductByIDHandler(t *testing.T) {
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
		router.HandleFunc(apiUrl, handler.GetProductByIDHandler).Methods(http.MethodGet)

		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Error(rr.Body)
			t.Errorf("expected status code %d, got %d", http.StatusNotFound, rr.Code)
		}
	})

	t.Run("should succeed because the uuid is correct", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, corrApiUrl, strings.NewReader(""))
		if err != nil {
			t.Fatal(err)
		}
		defer req.Body.Close()

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc(apiUrl, handler.GetProductByIDHandler).Methods(http.MethodGet)

		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Error(rr.Body)
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})
}

func TestDeleteProductHandler(t *testing.T) {
	wrong := "/delete/" + wrongUUID
	corr := "/delete/" + mockUUID.String()
	url := "/delete/{iD}"

	t.Run("should fail because id is not correct", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, wrong, strings.NewReader(""))
		if err != nil {
			t.Fatal(err)
		}
		defer req.Body.Close()

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc(url, handler.DeleteProductHandler).Methods(http.MethodDelete)

		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusInternalServerError {
			t.Error(rr.Body)
			t.Errorf("expected code %d, got %d", http.StatusInternalServerError, rr.Code)
		}
	})

	t.Run("should succeed because id is correct", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, corr, strings.NewReader(""))
		if err != nil {
			t.Fatal(err)
		}
		defer req.Body.Close()

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc(url, handler.DeleteProductHandler).Methods(http.MethodDelete)

		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Error(rr.Body)
			t.Errorf("expected code %d, got %d", http.StatusOK, rr.Code)
		}
	})
}

// mock product store
type mockProductStore struct{}

func (m *mockProductStore) CreateProduct(input c_product.Product) (pp_product.ProductDto, error) {
	dto := pp_product.ProductDto{}

	return dto, nil
}

func (m *mockProductStore) GetAllProducts() ([]pp_product.ProductDto, error) {
	dtos := []pp_product.ProductDto{}
	return dtos, nil
}

func (m *mockProductStore) GetProductByID(iD uuid.UUID) (pp_product.ProductDto, error) {
	dto := pp_product.ProductDto{}
	if iD != mockUUID {
		return dto, errors.New("not nil")
	}
	return dto, nil
}

func (m *mockProductStore) UpdateProduct(input c_product.Product) (pp_product.ProductDto, error) {
	dto := pp_product.ProductDto{}
	return dto, nil
}

func (m *mockProductStore) UpdateProductItemCount(iD uuid.UUID, c uint64) error {
	return nil
}

func (m *mockProductStore) DeleteProduct(iD uuid.UUID) error {
	if iD != mockUUID {
		return errors.New("not nil")
	}
	return nil
}

// mock admin store

type mockAdminStore struct{}

func (m *mockAdminStore) CreateAdmin(a c_admin.Admin) error {
	return nil
}

func (m *mockAdminStore) GetAdminByUsername(u string) (c_admin.Admin, error) {
	a := c_admin.Admin{}
	return a, nil
}

func (m *mockAdminStore) GetAdminByID(iD uuid.UUID) (c_admin.Admin, error) {
	a := c_admin.Admin{}
	return a, nil
}

func (m *mockAdminStore) UpdateAdminPassword(iD uuid.UUID, np string) error {
	return nil
}
