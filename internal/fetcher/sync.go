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

func printDuplicates(dishBindings []models.DishBinding) {
	mb := make(map[uuid.UUID]struct{})
	for _, db := range dishBindings {
		_, ok := mb[db.GlovoDishID]
		if ok && db.GlovoDishID != uuid.Nil {
			fmt.Printf("Duplicate key: %v\n", db.GlovoDishID)
		} else {
			mb[db.GlovoDishID] = struct{}{}
			fmt.Printf("Added key: %v\n", db.GlovoDishID)
		}
	}

}

func bindRestaurants(glovoRestauraunts []models.GlovoRestaurant, yandexRestaurants []models.YandexRestaurant) (overlap, glovoOnly, yandexOnly []models.RestaurantBinding) {
	mb := make(map[string]uuid.UUID)
	for _, grest := range glovoRestauraunts {
		mb[grest.Name] = grest.ID
	}

	extracted := make(map[string]struct{})
	for _, yrest := range yandexRestaurants {
		grestID, ok := mb[yrest.Name]
		if ok {
			overlap = append(overlap, models.RestaurantBinding{GlovoRestaurantID: grestID, YandexRestaurantID: yrest.ID})
			extracted[yrest.Name] = struct{}{}
		} else {
			yandexOnly = append(yandexOnly, models.RestaurantBinding{YandexRestaurantID: yrest.ID})
		}
	}

	for _, v := range glovoRestauraunts {
		_, ok := extracted[v.Name]
		if !ok {
			glovoOnly = append(glovoOnly, models.RestaurantBinding{GlovoRestaurantID: v.ID})
		}
	}

	return overlap, glovoOnly, yandexOnly

}

func bindDishes(rbID uuid.UUID, gdishes []models.GlovoDish, ydishes []models.YandexDish) (overlap, gOnly, yOnly []models.DishBinding) {
	mb := make(map[string]uuid.UUID)
	for _, gdish := range gdishes {
		mb[gdish.Name] = gdish.ID
	}

	extracted := make(map[uuid.UUID]struct{})
	for _, ydish := range ydishes {
		gdishID, ok := mb[ydish.Name]
		if ok {
			b, err := makeDishBinding(rbID, gdishID, ydish.ID)
			if err != nil {
				log.Printf("Error caught when creating dish binding: %v", err)
				continue
			}

			overlap = append(overlap, b)
			extracted[gdishID] = struct{}{}
		} else {
			b, err := makeDishBinding(rbID, uuid.Nil, ydish.ID)
			if err != nil {
				log.Printf("Error caught when creating dish binding: %v", err)
				continue
			}
			yOnly = append(yOnly, b)
		}
	}

	for _, v := range gdishes {
		_, ok := extracted[v.ID]
		if !ok {
			b, err := makeDishBinding(rbID, v.ID, uuid.Nil)
			if err != nil {
				log.Printf("Error caught when creating dish binding: %v", err)
				continue
			}
			gOnly = append(gOnly, b)
		}
	}

	return overlap, gOnly, yOnly
}

func makeDishBinding(rbID, gdishID, ydishID uuid.UUID) (models.DishBinding, error) {
	if rbID == uuid.Nil {
		return models.DishBinding{}, fmt.Errorf("Restaurant binding ID can't be null when creating dish binding")
	}

	if gdishID == uuid.Nil && ydishID == uuid.Nil {
		return models.DishBinding{}, fmt.Errorf("Both yandex dish ID and glovo dish ID can't be null when creating a dish binding")
	}

	return models.DishBinding{
		ID:                  uuid.New(),
		RestaurantBindingID: rbID,
		GlovoDishID:         gdishID,
		YandexDishID:        ydishID,
	}, nil
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
	restaurantBindings, err := cfg.DB.GetRestaurantBindingsToUpdate(ctx)
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
	fmt.Printf("CREATING DISHES FOR %v\n", uuid.UUID(rb.ID.Bytes))
	ctx := context.Background()
	glovoDishes := []models.GlovoDish{}
	yandexDishes := []models.YandexDish{}
	if rb.GlovoRestaurantID.Valid {
		fmt.Println("Glovo RestaurantID valid")
		dbDishes, err := cfg.DB.GetGlovoDishesForRestaurant(ctx, rb.GlovoRestaurantID)
		if err != nil {
			log.Fatal(err)
		}
		dishes := []models.GlovoDish{}
		for _, d := range dbDishes {
			dish := utils.GlovoDishDBtoModel(d)
			dishes = append(dishes, dish)
		}
		glovoDishes = append(glovoDishes, dishes...)
		// for _, d := range glovoDishes {
		// 	fmt.Printf("Glovo DishID: %v\n", d.ID)
		// }
	}

	if rb.YandexRestaurantID.Valid {
		fmt.Println("Yandex RestaurantID valid")
		dbDishes, err := cfg.DB.GetYandexDishesForRestaurant(ctx, rb.YandexRestaurantID)
		if err != nil {
			log.Fatal(err)
		}
		dishes := []models.YandexDish{}
		for _, d := range dbDishes {
			dish := utils.YandexDishDBtoModel(d)
			dishes = append(dishes, dish)
		}
		yandexDishes = append(yandexDishes, dishes...)
		// for _, d := range yandexDishes {
		// 	fmt.Printf("Yandex DishID: %v\n", d.ID)
		// }
	}

	ov, glo, yo := bindDishes(rb.ID.Bytes, glovoDishes, yandexDishes)
	bindings := []models.DishBinding{}
	bindings = append(bindings, ov...)
	bindings = append(bindings, glo...)
	bindings = append(bindings, yo...)
	// for i := 0; i < len(bindings); i++ {
	// 	b := &bindings[i]
	// 	b.ID = uuid.New()
	// 	b.RestaurantBindingID = rb.ID.Bytes
	// 	// fmt.Println(b)
	// }

	// for _, b := range bindings {
	// 	fmt.Printf("binding: %v\n", b)
	// }

	// printDuplicates(bindings)

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
