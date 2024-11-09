package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/thrashdev/foodsearch/internal/fetcher"
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
	dishes, err := fetcher.FetchGlovoDishes(restaurants[0], glovoDishesURL)
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

	for _, d := range dishes {
		fmt.Println(d)
	}
}
