package main

import (
	"log"
	"os"

	"github.com/o-rensa/iv/cmd/api"
	"github.com/o-rensa/iv/pkg/initializers"
)

func init() {
	initializers.LoadDotEnv()
}

func main() {
	server := api.NewAPIServer(os.Getenv("PORT"), nil)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
