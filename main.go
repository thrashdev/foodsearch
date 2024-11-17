package main

import (
	// "fmt"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/thrashdev/foodsearch/internal/config"
	"github.com/thrashdev/foodsearch/internal/database"
	"github.com/thrashdev/foodsearch/internal/fetcher"

	// "github.com/thrashdev/foodsearch/internal/fetcher"
	// "github.com/thrashdev/foodsearch/internal/models"
	"net/http"
	// "github.com/thrashdev/foodsearch/internal/models"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	type responseReady struct {
		Status string `json:"status"`
	}

	resp := responseReady{"OK"}
	err := respondWithJSON(w, 200, resp)
	if err != nil {
		log.Println("Server failed on checking readiness")
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	glovoSearchURL := os.Getenv("glovo_url")
	glovoFiltersURL := os.Getenv("glovo_filters_url")
	glovoDishURL := os.Getenv("glovo_dishes_url")
	port := os.Getenv("PORT")
	connection_string := os.Getenv("connection_string")
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connection_string)
	if err != nil {
		log.Fatalf("Couldn't connect to database :%v", err)
	}
	defer conn.Close(ctx)
	db := database.New(conn)
	config := config.Config{
		Glovo: config.GlovoConfig{
			SearchURL:  glovoSearchURL,
			FiltersURL: glovoFiltersURL,
			DishURL:    glovoDishURL,
		},
		DB:              *db,
		UpdateBatchSize: 5,
	}
	fmt.Println("Started fetching new restaurants")
	go fetcher.FetchNewGlovoRestaurants(config)
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("GET /v1/healthz", handlerReadiness)

	server := http.Server{Handler: serveMux, Addr: ":" + port}
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Couldn't start the server: ", err)
	}

}
