package ap_admin

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	ap_auth "github.com/o-rensa/iv/internal/adapters/primary/auth"
	ap_shared "github.com/o-rensa/iv/internal/adapters/primary/shared"
	c_admin "github.com/o-rensa/iv/internal/core/admin"
	c_sharedTypes "github.com/o-rensa/iv/internal/core/sharedtypes"
	pp_admin "github.com/o-rensa/iv/internal/ports/primary/admin"
	ps_admin "github.com/o-rensa/iv/internal/ports/secondary/admin"
	"github.com/o-rensa/iv/pkg/configs"
	"github.com/o-rensa/iv/pkg/middlewares"
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
	subrouter.HandleFunc("/updatepassword", ap_auth.WithJWTAuth(ah.ChangeAdminPasswordHandler, ah.adminStore)).Methods(http.MethodPost)
	subrouter.HandleFunc("/login", ah.HandleLogin).Methods(http.MethodPost)
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
	_, err := ah.adminStore.GetAdminByUsername(payload.Username)
	if err == nil {
		ap_shared.WriteError(w, http.StatusBadRequest, fmt.Errorf("admin with username %s already exists", payload.Username))
	}

	// hash password
	hashedPassword, err := middlewares.HashPassword(payload.UnhashedPassword)
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

func (ah *AdminHandler) ChangeAdminPasswordHandler(w http.ResponseWriter, r *http.Request) {
	// decode request and store into update password payload
	var pl pp_admin.UpdateAdminPasswordPayload
	if err := ap_shared.ParseJSON(r, &pl); err != nil {
		ap_shared.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload
	if err := ap_shared.Validate.Struct(pl); err != nil {
		errors := err.(validator.ValidationErrors)
		ap_shared.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	// convert string to uuid
	aID := uuid.MustParse(pl.ID)

	// get admin by ID
	a, err := ah.adminStore.GetAdminByID(aID)
	if err != nil {
		ap_shared.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate password
	if cp := middlewares.ComparePasswords(a.ID.String(), pl.UnhashedOldPassword); !cp {
		ap_shared.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid old password"))
		return
	}

	// hash new password
	hashedNP, err := middlewares.HashPassword(pl.NewUnhashedPassword)
	if err != nil {
		ap_shared.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to hash new password"))
	}

	// update to new password
	updateErr := ah.adminStore.UpdateAdminPassword(aID, hashedNP)
	if updateErr != nil {
		ap_shared.WriteError(w, http.StatusInternalServerError, updateErr)
	}

	ap_shared.WriteJSON(w, http.StatusAccepted, nil)
}

func (ah *AdminHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var admin pp_admin.LogInAdminPayload

	// parser request to payload
	if err := ap_shared.ParseJSON(r, &admin); err != nil {
		ap_shared.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload
	if err := ap_shared.Validate.Struct(admin); err != nil {
		errors := err.(validator.ValidationErrors)
		ap_shared.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// get admin by username
	a, err := ah.adminStore.GetAdminByUsername(admin.Username)
	if err != nil {
		ap_shared.WriteError(w, http.StatusBadRequest, fmt.Errorf("username or password not found"))
		return
	}

	// check password
	if !middlewares.ComparePasswords((*a).HashedPassword, admin.UnhashedPassword) {
		ap_shared.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid password"))
		return
	}

	// create a jwt token
	secret := []byte(configs.Configs.JWTSecret)
	token, err := ap_auth.CreateJWT(secret, a.ID)
	if err != nil {
		ap_shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	ap_shared.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}
