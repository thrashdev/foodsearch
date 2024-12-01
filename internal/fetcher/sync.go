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

func bindRestaurants(glovoRestauraunts []models.GlovoRestaurant, yandexRestaurants []models.YandexRestaurant) (overlap, glovoOnly, yandexOnly []models.RestaurantBinding) {
	mb := make(map[string]models.RestaurantBinding)
	for _, grest := range glovoRestauraunts {
		mb[grest.Name] = models.RestaurantBinding{GlovoRestaurantID: grest.ID}
	}

	for _, yrest := range yandexRestaurants {
		b, ok := mb[yrest.Name]
		if ok {
			overlap = append(overlap, models.RestaurantBinding{GlovoRestaurantID: b.GlovoRestaurantID, YandexRestaurantID: yrest.ID})
		} else {
			yandexOnly = append(yandexOnly, models.RestaurantBinding{YandexRestaurantID: yrest.ID})
		}
	}

	for _, v := range mb {
		glovoOnly = append(glovoOnly, v)
	}

	return overlap, glovoOnly, yandexOnly

}

func bindDishes(gdishes []models.GlovoDish, ydishes []models.YandexDish) (overlap, gOnly, yOnly []models.DishBinding) {
	mb := make(map[string]models.DishBinding)
	for _, gdish := range gdishes {
		mb[gdish.Name] = models.DishBinding{ID: gdish.ID}
	}

	for _, ydish := range ydishes {
		b, ok := mb[ydish.Name]
		if ok {
			overlap = append(overlap, models.DishBinding{GlovoDishID: b.GlovoDishID, YandexDishID: ydish.ID})
		} else {
			yOnly = append(yOnly, models.DishBinding{YandexDishID: ydish.ID})
		}
	}

	for _, v := range mb {
		gOnly = append(gOnly, v)
	}

	return overlap, gOnly, yOnly
}

func SyncRestaurants(cfg *config.Config) {
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

	ov, glo, yo := bindRestaurants(glovoRests, yandexRests)
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

func SyncDishes(cfg *config.Config) (rowsAffected int64) {
	ctx := context.Background()
	restaurantBindings, err := cfg.DB.GetAllRestaurantBindings(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, rb := range restaurantBindings {
		res := syncDishes(cfg, rb)
		rowsAffected += res
	}

	fmt.Printf("Created %v dishes\n", rowsAffected)
	return rowsAffected

}

func syncDishes(cfg *config.Config, rb database.RestaurantBinding) int64 {
	ctx := context.Background()
	glovoDishes := []models.GlovoDish{}
	yandexDishes := []models.YandexDish{}
	if rb.GlovoRestaurantID.Valid {
		dbDishes, err := cfg.DB.GetAllGlovoDishes(ctx)
		if err != nil {
			log.Fatal(err)
		}
		dishes := []models.GlovoDish{}
		for _, d := range dbDishes {
			dish := utils.GlovoDishDBtoModel(d)
			dishes = append(dishes, dish)
		}
		glovoDishes = append(glovoDishes, dishes...)
	}

	if rb.YandexRestaurantID.Valid {
		dbDishes, err := cfg.DB.GetAllYandexDishes(ctx)
		if err != nil {
			log.Fatal(err)
		}
		dishes := []models.YandexDish{}
		for _, d := range dbDishes {
			dish := utils.YandexDishDBtoModel(d)
			dishes = append(dishes, dish)
		}
		yandexDishes = append(yandexDishes, dishes...)
	}

	ov, glo, yo := bindDishes(glovoDishes, yandexDishes)
	bindings := []models.DishBinding{}
	bindings = append(bindings, ov...)
	bindings = append(bindings, glo...)
	bindings = append(bindings, yo...)
	for i := 0; i < len(bindings); i++ {
		b := &bindings[i]
		b.ID = uuid.New()
		b.RestaurantBindingID = rb.ID.Bytes
		fmt.Println(b)
	}

	args := []database.BatchCreateDishBindingsParams{}
	for _, b := range bindings {
		args = append(args, utils.DishBindingsModelToDB(b))
	}
	rowsAffected, err := cfg.DB.BatchCreateDishBindings(ctx, args)
	if err != nil {
		log.Fatal(err)
	}

	return rowsAffected
}
