package tests

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	ap_admin "github.com/o-rensa/iv/internal/adapters/primary/admin"
	c_admin "github.com/o-rensa/iv/internal/core/admin"
)

func TestCreateAdminServiceHandler(t *testing.T) {
	store := &mockAdminStore{}
	handler := ap_admin.NewAdminHandler(store)
	apiUrl := "/create"

	noUsernameCase := []byte(`{"password": "qwerty"}`)
	noPasswordCase := []byte(`{"username": "admin"}`)
	validCase := []byte(`{"username": "admin", "password": "qwerty"}`)

	t.Run("should fail because username is missing", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(noUsernameCase))
		if err != nil {
			t.Fatal(err)
		}
		defer req.Body.Close()

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc(apiUrl, handler.CreateAdminHandler).Methods(http.MethodPost)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Error(rr.Body)
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should fail because password is missing", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(noPasswordCase))
		if err != nil {
			t.Fatal(err)
		}
		defer req.Body.Close()

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc(apiUrl, handler.CreateAdminHandler).Methods(http.MethodPost)

		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Error(rr.Body)
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should not fail cause payload is valid", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(validCase))
		if err != nil {
			t.Fatal(err)
		}
		defer req.Body.Close()

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc(apiUrl, handler.CreateAdminHandler).Methods(http.MethodPost)

		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusCreated {
			t.Error(rr.Body)
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})

}

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
