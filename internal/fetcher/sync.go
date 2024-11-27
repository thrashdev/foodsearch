package fetcher

import (
	"context"
	"log"

	"github.com/thrashdev/foodsearch/internal/config"
	"github.com/thrashdev/foodsearch/internal/models"
	"github.com/thrashdev/foodsearch/internal/utils"
)

func getSubset(glovoRestauraunts []models.GlovoRestaurant, yandexRestaurants []models.YandexRestaurant) (overlap, glovoOnly, yandexOnly []models.RestaurantBinding) {
	mb := make(map[string]models.RestaurantBinding)
	for _, grest := range glovoRestauraunts {
		mb[grest.Name] = models.RestaurantBinding{GlovoRestaurantID: grest.ID}
	}

	for _, yrest := range yandexRestaurants {
		b, ok := mb[yrest.Name]
		if ok {
			overlap = append(overlap, models.RestaurantBinding{GlovoRestaurantID: b.GlovoRestaurantID, YandexRestaurantID: yrest.ID})
		} else {
			glovoOnly = append(glovoOnly, models.RestaurantBinding{GlovoRestaurantID: b.GlovoRestaurantID})
			yandexOnly = append(yandexOnly, models.RestaurantBinding{YandexRestaurantID: yrest.ID})
		}
	}

	return overlap, glovoOnly, yandexOnly

}

func sync(cfg *config.Config) {
	ctx := context.Background()
	glovoRestaurants, err := cfg.DB.GetAllGlovoRestaurants(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	yandexRestaurants, err := cfg.DB.GetAllYandexRestaurants(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	glovoRests := []models.GlovoRestaurant{}
	for _, g := range glovoRestaurants {
		glovoRests = append(glovoRests, utils.GlovoRestDBtoModel(g))
	}

	yandexRests := []models.YandexRestaurant{}
	for _, y := range yandexRestaurants {
		yandexRests = append(yandexRests, utils.YandexRestDBtoModel(y))
	}

	ov, glo, yo := getSubset(glovoRests, yandexRests)
	bindings := []models.RestaurantBinding{}
	bindings = append(bindings, ov...)
	bindings = append(bindings, glo...)
	bindings = append(bindings, yo...)

}
