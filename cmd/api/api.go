package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	ap_admin "github.com/o-rensa/iv/internal/adapters/primary/admin"
	ap_brand "github.com/o-rensa/iv/internal/adapters/primary/brands"
	as_admin "github.com/o-rensa/iv/internal/adapters/secondary/admin"
	as_brand "github.com/o-rensa/iv/internal/adapters/secondary/brands"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: "localhost:" + addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	//admin
	adminStore, _ := as_admin.NewAdminStore(s.db)
	adminHandler := ap_admin.NewAdminHandler(adminStore)
	adminHandler.RegisterRoutes(subrouter)

	//brand
	brandStore, _ := as_brand.NewBrandStore(s.db)
	brandhandler := ap_brand.NewBrandHandler(brandStore, adminStore)
	brandhandler.RegisterRoutes(subrouter)

	log.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, router)
}
