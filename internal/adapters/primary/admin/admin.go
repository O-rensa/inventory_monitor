package ap_admin

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	ap_auth "github.com/o-rensa/iv/internal/adapters/primary/auth"
	ap_shared "github.com/o-rensa/iv/internal/adapters/primary/shared"
	c_admin "github.com/o-rensa/iv/internal/core/admin"
	c_sharedTypes "github.com/o-rensa/iv/internal/core/sharedtypes"
	pp_admin "github.com/o-rensa/iv/internal/ports/primary/admin"
	ps_admin "github.com/o-rensa/iv/internal/ports/secondary/admin"
)

type AdminHandler struct {
	adminStore ps_admin.AdminStore
}

func NewAdminHandler(adminStore ps_admin.AdminStore) *AdminHandler {
	adminHandler := &AdminHandler{adminStore: adminStore}
	return adminHandler
}

func (ah *AdminHandler) RegisterRoutes(router *mux.Router) {
	subrouter := router.PathPrefix("/admin").Subrouter()

	subrouter.HandleFunc("/create", ah.CreateAdminHandler).Methods(http.MethodPost)
}

func (ah *AdminHandler) CreateAdminHandler(w http.ResponseWriter, r *http.Request) {
	// decode request and store into create admin payload
	var payload pp_admin.CreateAdminPayload
	if err := ap_shared.ParseJSON(r, &payload); err != nil {
		ap_shared.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload
	if err := ap_shared.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		ap_shared.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	// check if username exists
	err := ah.adminStore.GetAdminByUsername(payload.Username)
	if err == nil {
		ap_shared.WriteError(w, http.StatusBadRequest, fmt.Errorf("admin with username %s already exists", payload.Username))
	}

	// hash password
	hashedPassword, err := ap_auth.HashPassword(payload.UnhashedPassword)
	if err != nil {
		ap_shared.WriteError(w, http.StatusInternalServerError, err)
	}

	// create admin
	admin := &c_admin.Admin{
		ModelID:        c_sharedTypes.NewModelID(),
		Username:       payload.Username,
		HashedPassword: hashedPassword,
	}
	if err := ah.adminStore.CreateAdmin(admin); err != nil {
		ap_shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	ap_shared.WriteJSON(w, http.StatusCreated, nil)
}
