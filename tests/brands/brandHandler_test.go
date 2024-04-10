package tests

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	ap_brand "github.com/o-rensa/iv/internal/adapters/primary/brands"
	c_admin "github.com/o-rensa/iv/internal/core/admin"
	c_brands "github.com/o-rensa/iv/internal/core/brands"
	pp_brand "github.com/o-rensa/iv/internal/ports/primary/brands"
)

var brandStore = &mockBrandStore{}
var adminStore = &mockAdminStore{}
var handler = ap_brand.NewBrandHandler(brandStore, adminStore)
var mockUUID = uuid.MustParse("56bd6b24-83e5-4da7-ab26-e45cbe0aa8b6")
var wrongUUID = "7e742f61-58ed-4072-9980-b4170176ca44"

func TestCreateBrandHandler(t *testing.T) {
	apiUrl := "/create"

	// cases
	noBrandNameCase := []byte(`{"brandName":}`)
	validCase := []byte(`{"brandName":"SuperBrand"}`)

	t.Run("should fail because brandName is missing", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(noBrandNameCase))
		if err != nil {
			t.Fatal(err)
		}
		defer req.Body.Close()

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc(apiUrl, handler.CreateBrandHandler).Methods(http.MethodPost)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Error(rr.Body)
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should succeed case payload is valid", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(validCase))
		if err != nil {
			t.Fatal(err)
		}
		defer req.Body.Close()

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc(apiUrl, handler.CreateBrandHandler).Methods(http.MethodPost)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Error(rr.Body)
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})
}

func TestGetAllBrandsHandler(t *testing.T) {
	apiUrl := "/getall"

	t.Run("should succeed", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, apiUrl, strings.NewReader(""))
		if err != nil {
			log.Fatal(err)
		}
		defer req.Body.Close()

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc(apiUrl, handler.GetAllBrandsHandler).Methods(http.MethodGet)

		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("%v", rr.Body)
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})
}

func TestGetBrandByIDHandler(t *testing.T) {
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
		router.HandleFunc(apiUrl, handler.GetBrandByIDHandler).Methods(http.MethodGet)

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
		router.HandleFunc(apiUrl, handler.GetBrandByIDHandler).Methods(http.MethodGet)

		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("%v", rr.Body)
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})
}

func TestUpdateBrandHandler(t *testing.T) {
	wrongPayload := []byte(`{"iD":"` + wrongUUID + `", "brandName":"New Brand Name"}`)
	corrPayload := []byte(`{"iD":"` + mockUUID.String() + `", "brandName":"New Brand Name"}`)

	apiUrl := "/update"

	t.Run("should fail because the ID is not correct", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPut, apiUrl, bytes.NewBuffer(wrongPayload))
		if err != nil {
			t.Fatal(err)
		}
		defer req.Body.Close()
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc(apiUrl, handler.UpdateBrandHandler).Methods(http.MethodPut)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Error(rr.Body)
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
		}
	})

	t.Run("should succeed because ID is correct", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPut, apiUrl, bytes.NewBuffer(corrPayload))
		if err != nil {
			t.Fatal(err)
		}
		defer req.Body.Close()

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc(apiUrl, handler.UpdateBrandHandler).Methods(http.MethodPut)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Error(rr.Body)
			t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
		}
	})
}

func TestDeleteBrandHandler(t *testing.T) {
	wrongApiUrl := "/delete/" + wrongUUID
	corrApiUrl := "/delete/" + mockUUID.String()
	apiUrl := "/delete/{iD}"

	t.Run("should fail because id is not correct", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, wrongApiUrl, strings.NewReader(""))
		if err != nil {
			t.Fatal(err)
		}
		defer req.Body.Close()

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc(apiUrl, handler.DeleteBrandHandler).Methods(http.MethodDelete)

		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusInternalServerError {
			t.Error(rr.Body)
			t.Errorf("expected code %d, got %d", http.StatusInternalServerError, rr.Code)
		}
	})

	t.Run("should succeed because if is correct", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, corrApiUrl, strings.NewReader(""))
		if err != nil {
			t.Fatal(err)
		}
		defer req.Body.Close()

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc(apiUrl, handler.DeleteBrandHandler).Methods(http.MethodDelete)

		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Error(rr.Body)
			t.Errorf("expected code %d, got %d", http.StatusOK, rr.Code)
		}
	})
}

// mock brand store
type mockBrandStore struct{}

func (m *mockBrandStore) CreateBrand(i c_brands.Brand) (pp_brand.BrandDto, error) {
	dto := pp_brand.BrandDto{}
	dto.BrandID = i.ModelID.ID.String()
	dto.BrandName = i.BrandName
	return dto, nil
}

func (m *mockBrandStore) GetAllBrands() ([]pp_brand.BrandDto, error) {
	dtos := []pp_brand.BrandDto{}
	return dtos, nil
}

func (m *mockBrandStore) GetBrandByID(i uuid.UUID) (pp_brand.BrandDto, error) {
	dto := pp_brand.BrandDto{}
	if i != mockUUID {
		return dto, errors.New("id did not match to mockUUID")
	}
	dto.BrandID = mockUUID.String()
	dto.BrandName = "brandUno"
	return dto, nil
}

func (m *mockBrandStore) UpdateBrand(i c_brands.Brand) (pp_brand.BrandDto, error) {
	dto := pp_brand.BrandDto{}
	if i.ModelID.ID != mockUUID {
		return dto, errors.New("not nil error")
	}
	dto.BrandID = i.ModelID.ID.String()
	dto.BrandName = i.BrandName
	return dto, nil
}

func (m *mockBrandStore) DeleteBrand(i uuid.UUID) error {
	if i != mockUUID {
		return errors.New("not nil")
	}
	return nil
}

// mock admin store
type mockAdminStore struct{}

func (m *mockAdminStore) CreateAdmin(a c_admin.Admin) error {
	return nil
}

func (m *mockAdminStore) GetAdminByUsername(username string) (c_admin.Admin, error) {
	av := c_admin.Admin{}
	return av, errors.New("error")
}

func (m *mockAdminStore) GetAdminByID(Id uuid.UUID) (c_admin.Admin, error) {
	av := c_admin.Admin{}
	return av, nil
}

func (m *mockAdminStore) UpdateAdminPassword(Id uuid.UUID, np string) error {
	return nil
}
