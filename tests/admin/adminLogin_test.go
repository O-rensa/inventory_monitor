package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	ap_admin "github.com/o-rensa/iv/internal/adapters/primary/admin"
	c_admin "github.com/o-rensa/iv/internal/core/admin"
	c_sharedTypes "github.com/o-rensa/iv/internal/core/sharedtypes"
	"github.com/o-rensa/iv/pkg/middlewares"
)

func TestLogInAdminServiceHandler(t *testing.T) {
	store := &mockAdminLoginStore{}
	handler := ap_admin.NewAdminHandler(store)
	apiUrl := "/login"

	noUsernameCase := []byte(`{"password": "qwerty"}`)
	wrongPasswordCase := []byte(`{"username": "admin", "password": "notqwerty"}`)
	validCase := []byte(`{"username": "admin", "password": "qwerty"}`)

	t.Run("should fail because payload has no username", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(noUsernameCase))
		if err != nil {
			t.Fatal(err)
		}
		defer req.Body.Close()

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc(apiUrl, handler.HandleLogin).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Error(rr.Body)
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should fail because password is wrong", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(wrongPasswordCase))
		if err != nil {
			t.Fatal(err)
		}
		defer req.Body.Close()

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc(apiUrl, handler.HandleLogin).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Error(rr.Body)
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should not fail because payload is valid", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(validCase))
		if err != nil {
			t.Fatal(err)
		}
		defer req.Body.Close()

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc(apiUrl, handler.HandleLogin).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Error(rr.Body)
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})
}

type mockAdminLoginStore struct{}

func (m *mockAdminLoginStore) CreateAdmin(a c_admin.Admin) error {
	return nil
}

func (m *mockAdminLoginStore) GetAdminByUsername(username string) (c_admin.Admin, error) {
	hPasswd, _ := middlewares.HashPassword("qwerty")
	av := c_admin.Admin{
		ModelID:        c_sharedTypes.NewModelID(),
		Username:       "admin",
		HashedPassword: hPasswd,
	}
	return av, nil
}

func (m *mockAdminLoginStore) GetAdminByID(Id uuid.UUID) (c_admin.Admin, error) {
	av := c_admin.Admin{}
	return av, nil
}

func (m *mockAdminLoginStore) UpdateAdminPassword(Id uuid.UUID, np string) error {
	return nil
}
