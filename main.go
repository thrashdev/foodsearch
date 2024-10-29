package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/thrashdev/foodsearch/internal/fetcher"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	glovoURL := os.Getenv("glovo_url")
	glovoFiltersURL := os.Getenv("glovo_filters_url")
	restaurants, err := fetcher.FetchGlovoItems(glovoURL, glovoFiltersURL)
	if err != nil {
		fmt.Println(err)
	}
	for _, r := range restaurants {
		fmt.Println(r)
	}
}
