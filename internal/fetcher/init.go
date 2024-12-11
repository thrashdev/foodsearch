package fetcher

import "github.com/thrashdev/foodsearch/internal/config"

// type startupFunction func(cfg *config.Config) (

// TODO: catch all returns in this function and print it here
func Init(cfg *config.Config) {
	InitGlovo(cfg)
	InitYandex(cfg)
	SyncRestaurants(cfg)
	SyncDishes(cfg)
}
