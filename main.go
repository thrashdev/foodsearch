package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/thrashdev/foodsearch/internal/fetcher"
	"github.com/thrashdev/foodsearch/internal/models"
	"net/http"
	// "github.com/thrashdev/foodsearch/internal/models"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	port := os.Getenv("PORT")
	serveMux := http.NewServeMux()
	server := http.Server{Handler: serveMux, Addr: ":" + port}
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Couldn't start the server: ", err)
	}
	// glovoURL := os.Getenv("glovo_url")
	// glovoFiltersURL := os.Getenv("glovo_filters_url")
	// glovoDishesURL := os.Getenv("glovo_dishes_url")

}
