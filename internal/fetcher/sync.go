package fetcher

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/thrashdev/foodsearch/internal/config"
	"github.com/thrashdev/foodsearch/internal/database"
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

func Sync(cfg *config.Config) {
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
	for i := 0; i < len(bindings); i++ {
		b := &bindings[i]
		b.ID = uuid.New()
		fmt.Println(b)
	}

	fmt.Println("Overlap: ")
	fmt.Println(ov)

	args := []database.BatchCreateRestaurantBindingParams{}
	for _, b := range bindings {
		arg := utils.RestaurantBindingModeltoDB(b)
		args = append(args, arg)
	}

	rowsAffected, err := cfg.DB.BatchCreateRestaurantBinding(context.Background(), args)
	if err != nil {
		log.Fatalf("Couldn't create restaurant bindings in DB: %v", err)
	}
	fmt.Printf("Created %v restaurant bindings", rowsAffected)

}
