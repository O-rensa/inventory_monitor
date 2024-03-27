package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/o-rensa/iv/cmd/api"
	"github.com/o-rensa/iv/pkg/initializers"
)

var dB *sql.DB
var dBErr error

func init() {
	// Load environment variables
	initializers.LoadDotEnv()

	// initialize database
	dB, dBErr = initializers.InitializePostGres()
	if dBErr != nil {
		log.Fatal(dBErr)
	}
}

func main() {
	server := api.NewAPIServer(os.Getenv("PORT"), dB)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
