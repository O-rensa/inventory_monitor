package tests

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	ap_admin "github.com/o-rensa/iv/internal/adapters/primary/admin"
)

func TestLogInAdminServiceHandler(t *testing.T) {
	store := &mockAdminStore{}
	handler := ap_admin.NewAdminHandler(store)
	apiUrl := "/login"

	wrongPayloadCase := []byte(`"{username: admin, password: notqwerty}"`)
	validCase := []byte(`"{username: admin, password: qwerty}"`)

	t.Run("should fail because payload is invalid", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(wrongPayloadCase))
		if err != nil {
			log.Fatal(err)
		}
		defer req.Body.Close()

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc(apiUrl, handler.HandleLogin).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should not fail because payload is valid", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(validCase))
		if err != nil {
			log.Fatal(err)
		}
		defer req.Body.Close()

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc(apiUrl, handler.HandleLogin).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})
}
