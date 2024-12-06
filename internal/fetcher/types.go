package fetcher

//
// import (
// 	"fmt"
//
// 	"github.com/thrashdev/foodsearch/internal/config"
// 	"github.com/thrashdev/foodsearch/internal/utils"
// )
//
// type DBActionResult struct {
// 	action       string
// 	rowsAffected int64
// }
//
// func Init(cfg *config.Config) {
// 	var rowsAffected int64
// 	errCh := make(chan error)
// 	go utils.PrintErrors(errCh)
// 	CreateNewGlovoRestaurants(cfg, errCh)
// 	CreateNewDishesForGlovoRestaurants(cfg)
// 	fmt.Println("Started fetching new restaurants")
// 	rowsAffected = CreateNewYandexRestaurants(cfg)
// 	if err != nil {
// 		log.Fatalf("Error fetching yandex restaurants: %v", err)
// 	}
// 	fmt.Printf("Created %v restaurants\n", rowsAffected)
//
// 	rowsAffected = CreateNewYandexDishes(cfg)
// 	fmt.Printf("Created %v dishes\n", rowsAffected)
//
// 	fetcher.SyncRestaurants(cfg)
// 	fetcher.SyncDishes(cfg)
// }
