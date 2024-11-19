package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/thrashdev/foodsearch/internal/config"
	"github.com/thrashdev/foodsearch/internal/database"
	"github.com/thrashdev/foodsearch/internal/models"
	"github.com/thrashdev/foodsearch/internal/utils"
)

func restaurantDifference(restaurants []models.GlovoRestaurant, dbrNames []string) []models.GlovoRestaurant {
	mb := make(map[string]struct{}, len(dbrNames))
	for _, name := range dbrNames {
		mb[name] = struct{}{}
	}
	var diff []models.GlovoRestaurant
	for _, rest := range restaurants {
		if _, found := mb[rest.Name]; !found {
			diff = append(diff, rest)
		}
	}
	return diff

}

func dishDifference(dishes []models.GlovoDish, dbDishNames []string) []models.GlovoDish {
	mb := make(map[string]struct{}, len(dbDishNames))
	for _, name := range dbDishNames {
		mb[name] = struct{}{}
	}
	var diff []models.GlovoDish
	for _, dish := range dishes {
		if _, found := mb[dish.Name]; !found {
			diff = append(diff, dish)
		}
	}
	return diff
}

func findMaxDiscountRate(promos []glovoPromotion) float64 {
	max := 0.0
	for _, promo := range promos {
		if promo.Percentage > max {
			max = promo.Percentage
		}
	}
	// fmt.Printf("Returning %v\n", max)
	return max
}

func fetchByUrl(url string) (payload []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, fmt.Errorf("Couldn't create request, err: %v", err)
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Glovo-API-Version", "14")
	req.Header.Set("Glovo-App-Platform", "web")
	req.Header.Set("Glovo-App-Type", "customer")
	req.Header.Set("Glovo-App-Version", "7")
	req.Header.Set("Glovo-Location-City-Code", "BSK")
	req.Header.Set("Connection", "keep-alive")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, fmt.Errorf("Couldn't post request, err: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Couldn't read response body, err: %v", err)
	}

	return respBody, nil
}

func fetchGlovoFilters(filtersURL string) (filters []string, err error) {
	payload, err := fetchByUrl(filtersURL)
	if err != nil {
		return []string{}, err
	}
	var glovoFilters glovoFiltersResponse
	err = json.Unmarshal(payload, &glovoFilters)
	if err != nil {
		return []string{}, fmt.Errorf("Couldn't unmarshal json, err: %v", err)
	}

	for _, item := range glovoFilters.TopFilters {
		filters = append(filters, item.FilterName)
	}

	return filters, nil
}

func fetchGlovoRestaurantsByFilter(baseURL string, filter string) (restaurants []models.GlovoRestaurant, err error) {
	fullURL := baseURL + "&filter=" + filter
	respBody, err := fetchByUrl(fullURL)

	var glovoResp glovoRestaurantsResponse
	err = json.Unmarshal(respBody, &glovoResp)
	if err != nil {
		log.Println(fullURL)
		return []models.GlovoRestaurant{}, err
	}

	for _, item := range glovoResp.Elements {
		glovoRest := models.GlovoRestaurant{
			Restaurant: models.Restaurant{
				ID:          uuid.New(),
				Name:        item.SingleData.StoreData.Store.Name,
				Address:     item.SingleData.StoreData.Store.Address,
				DeliveryFee: item.SingleData.StoreData.Store.DeliveryFeeInfo.Fee,
				PhoneNumber: item.SingleData.StoreData.Store.PhoneNumber,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			GlovoApiStoreID:   item.SingleData.StoreData.Store.ID,
			GlovoApiAddressID: item.SingleData.StoreData.Store.AddressID,
			GlovoApiSlug:      item.SingleData.StoreData.Store.Slug,
		}

		restaurants = append(restaurants, glovoRest)
	}

	return restaurants, nil

}

// TODO: implement proper error-handling with an error channel
func CreateNewDishesForRestaurants(cfg *config.Config) error {
	ctx := context.Background()
	totalDishesCreated := 0
	dbRestaurants, err := cfg.DB.GetAllGlovoRestaurants(ctx)
	if err != nil {
		return err
	}
	for _, dbRest := range dbRestaurants {
		rest := utils.DatabaseGlovoRestaurantToModel(dbRest)
		dishesCreated := createNewDishesForGlovoRestaurant(cfg, rest)
		totalDishesCreated += dishesCreated
	}
	fmt.Printf("Created %v total dishes for %v restaurants", totalDishesCreated, len(dbRestaurants))
	return nil
}

// TODO: implement proper error-handling with an error channel
func createNewDishesForGlovoRestaurant(cfg *config.Config, rest models.GlovoRestaurant) (dishesCreated int) {
	newDishes, err := FetchGlovoDishes(rest, cfg.Glovo.DishURL)
	if err != nil {
		log.Printf("Couldn't fetch dishes: %v", err)
		return 0
	}
	ctx := context.Background()
	dbDishNames, err := cfg.DB.GetGlovoDishNames(ctx)
	dishesToAdd := dishDifference(newDishes, dbDishNames)
	args := []database.BatchCreateGlovoDishesParams{}
	for _, dish := range dishesToAdd {
		arg := database.BatchCreateGlovoDishesParams{
			ID:                pgtype.UUID{Bytes: uuid.New(), Valid: true},
			Name:              dish.Name,
			Description:       dish.Description,
			Price:             utils.FloatToNumeric(dish.Price),
			Discount:          utils.FloatToNumeric(dish.MaxDiscount),
			GlovoApiDishID:    int32(dish.GlovoAPIDishID),
			GlovoRestaurantID: pgtype.UUID{Bytes: dish.GlovoRestaurantID, Valid: true},
			CreatedAt:         pgtype.Timestamp{Time: time.Now().UTC(), InfinityModifier: 0, Valid: true},
			UpdatedAt:         pgtype.Timestamp{Time: time.Now().UTC(), InfinityModifier: 0, Valid: true},
		}

		args = append(args, arg)
	}

	rowsAffected, err := cfg.DB.BatchCreateGlovoDishes(ctx, args)
	if err != nil {
		log.Printf("Couldn't create dishes: %v", err)
		return 0
	}
	log.Printf("Created %v dishes", rowsAffected)
	return int(rowsAffected)

}

// TODO: implement proper error-handling with an error channel
func FetchNewGlovoRestaurants(cfg *config.Config) error {
	newRestaurants, err := fetchGlovoRestaurants(cfg.Glovo.SearchURL, cfg.Glovo.FiltersURL)
	if err != nil {
		return err
	}
	log.Printf("Fetched %v restaurants from API\n", len(newRestaurants))
	ctx := context.Background()
	dbRestaurantNames, err := cfg.DB.GetGlovoRestaurantNames(ctx)
	if err != nil {
		return err
	}
	log.Printf("Fetched %v restaurants from DB\n", len(dbRestaurantNames))
	restaurantsToAdd := restaurantDifference(newRestaurants, dbRestaurantNames)
	log.Printf("Restaurants to add: %v\n", len(restaurantsToAdd))
	args := []database.BatchCreateGlovoRestaurantsParams{}
	for _, rest := range restaurantsToAdd {
		arg := database.BatchCreateGlovoRestaurantsParams{
			ID:                pgtype.UUID{Bytes: uuid.New(), Valid: true},
			Name:              rest.Name,
			Address:           rest.Address,
			DeliveryFee:       utils.FloatToNumeric(rest.DeliveryFee),
			PhoneNumber:       pgtype.Text{String: rest.PhoneNumber, Valid: true},
			GlovoApiStoreID:   int32(rest.GlovoApiStoreID),
			GlovoApiAddressID: int32(rest.GlovoApiAddressID),
			GlovoApiSlug:      rest.GlovoApiSlug,
			CreatedAt:         pgtype.Timestamp{Time: time.Now().UTC(), InfinityModifier: 0, Valid: true},
			UpdatedAt:         pgtype.Timestamp{Time: time.Now().UTC(), InfinityModifier: 0, Valid: true},
		}

		args = append(args, arg)

		// res, err := cfg.DB.CreateGlovoRestaurant(ctx, arg)
		// if err != nil {
		// 	log.Println("Couldn't post restaurant to DB :v", err)
		// 	return err
		// }
		// newRows := res.RowsAffected()
		// rowsAffected += newRows
		// log.Printf("Created restaurant %v\n", rest.Name)
	}
	log.Printf("Prepared to create %v restaurants\n", len(args))
	rowsAffected, err := cfg.DB.BatchCreateGlovoRestaurants(ctx, args)
	if err != nil {
		fmt.Println(fmt.Errorf("Couldn't create glovo restaurants: %w", err))
	}

	log.Printf("Created %v restaurants\n", rowsAffected)
	return nil
}

// func UpdateRestaurants(cfg config.Config) error {
// 	ctx := context.Background()
// 	restaurantsToUpdate, err := cfg.DB.GetGlovoRestaurantsToUpdate(ctx, int32(cfg.UpdateBatchSize))
// 	if err != nil {
// 		return fmt.Errorf("Couldn't fetch glovo restaurants to update :w", err)
// 	}
//
// }
//
// func updateRestaurant(rest models.GlovoRestaurant) error {
//
// }

func fetchGlovoRestaurants(searchURL string, filtersURL string) (allRestaurants []models.GlovoRestaurant, err error) {
	filters, err := fetchGlovoFilters(filtersURL)
	if err != nil {
		return []models.GlovoRestaurant{}, fmt.Errorf("Couldn't get filters, err: %v", err)
	}

	for _, f := range filters {
		restaurantsByFilter, err := fetchGlovoRestaurantsByFilter(searchURL, url.QueryEscape(f))
		if err != nil {
			return []models.GlovoRestaurant{}, fmt.Errorf("Couldn't fetch by filter: %s. Error :%v", f, err)
		}

		allRestaurants = append(allRestaurants, restaurantsByFilter...)
	}
	return allRestaurants, nil
}

func FetchGlovoDishes(rest models.GlovoRestaurant, dishURL string) ([]models.GlovoDish, error) {
	targetURL := strings.Replace(dishURL, "{glovo_store_id}", strconv.Itoa(rest.GlovoApiStoreID), 1)
	targetURL = strings.Replace(targetURL, "{glovo_address_id}", strconv.Itoa(rest.GlovoApiAddressID), 1)
	responsePayload, err := fetchByUrl(targetURL)
	if err != nil {
		return []models.GlovoDish{}, fmt.Errorf("Error encountered while fetching glovo dishes: %v\n", err)
	}

	var dishesResponse glovoDishesResponse
	err = json.Unmarshal(responsePayload, &dishesResponse)
	if err != nil {
		return []models.GlovoDish{}, fmt.Errorf("Error encountered while fetching glovo dishes: %v\n", err)
	}

	dishes := []models.GlovoDish{}
	for _, elem := range dishesResponse.Data.Body {
		for _, dishItem := range elem.Data.Elements {
			discount := findMaxDiscountRate(dishItem.Data.Promotions)
			dishes = append(dishes, models.GlovoDish{
				GlovoAPIDishID: int(dishItem.Data.ID),
				Dish: models.Dish{
					ID:          uuid.New(),
					Name:        dishItem.Data.Name,
					Description: dishItem.Data.Description,
					Price:       dishItem.Data.Price,
					MaxDiscount: discount,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
				GlovoRestaurantID: rest.ID,
			})

		}
	}

	return dishes, nil

}

// func GlovoRespToGlovoStores(resp glovoRestaurantsResponse) ([]models.GlovoRestaurant, error) {
// 	result := []models.GlovoRestaurant{}
// 	for _, item := range resp.Elements {
// 		restaurant := models.GlovoRestaurant{
// 			Restaurant: models.Restaurant{
// 				ID:          uuid.New(),
// 				Name:        item.SingleData.StoreData.Store.Name,
// 				DeliveryFee: item.SingleData.StoreData.Store.DeliveryFeeInfo.Fee,
// 				Address:     item.SingleData.StoreData.Store.Address},
// 			GlovoStoreID:   item.SingleData.StoreData.Store.ID,
// 			GlovoAddressID: item.SingleData.StoreData.Store.AddressID,
// 		}
//
// 		result = append(result, restaurant)
// 	}
//
// 	return result, nil
// }
