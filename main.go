package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/thrashdev/foodsearch/internal/fetcher"
	"github.com/thrashdev/foodsearch/internal/models"
	// "github.com/thrashdev/foodsearch/internal/models"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	glovoURL := os.Getenv("glovo_url")
	glovoFiltersURL := os.Getenv("glovo_filters_url")
	glovoDishesURL := os.Getenv("glovo_dishes_url")
	restaurants, err := fetcher.FetchGlovoRestaurants(glovoURL, glovoFiltersURL)
	if err != nil {
		fmt.Println(err)
	}
	rest := models.GlovoRestaurant{}
	for _, r := range restaurants {
		if r.Name == "WoK Lagman" {
			rest = r
		}
	}
	dishes, err := fetcher.FetchGlovoDishes(rest, glovoDishesURL)
	if err != nil {
		fmt.Println(err)
	}

	// for _, r := range restaurants {
	// 	restDishes, err := fetcher.FetchGlovoDishes(r, glovoDishesURL)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		continue
	// 	}
	// 	dishes = append(dishes, restDishes...)
	// }

	for _, _ = range dishes {

	}
}
