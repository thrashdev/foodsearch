package fetcher

import (
	"context"
	"database/sql"
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
	"github.com/thrashdev/foodsearch/internal/config"
	"github.com/thrashdev/foodsearch/internal/database"
	"github.com/thrashdev/foodsearch/internal/models"
)

func findMaxDiscountRate(promos []glovoPromotion) float64 {
	max := 0.0
	for _, promo := range promos {
		if promo.Percentage > max {
			fmt.Printf("Processing %v\n", promo)
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
		}

		restaurants = append(restaurants, glovoRest)
	}

	return restaurants, nil

}

func FetchNewGlovoRestaurants(cfg config.Config) error {
	newRestaurants, err := fetchGlovoRestaurants(cfg.Glovo.SearchURL, cfg.Glovo.FiltersURL)
	if err != nil {
		return err
	}
	log.Println("Fetched restaurants from API")
	log.Println(newRestaurants)
	ctx := context.Background()
	dbRestaurants, err := cfg.DB.GetGlovoRestaurantNames(ctx)
	if err != nil {
		return err
	}
	log.Println("Fetched restaurant names from DB")
	restaurantsToAdd := []models.GlovoRestaurant{}
	for _, nr := range newRestaurants {
		for _, dbrName := range dbRestaurants {
			if nr.Name != dbrName {
				restaurantsToAdd = append(restaurantsToAdd, nr)
			}
		}
	}
	rowsAffected := int64(0)
	for _, rest := range restaurantsToAdd {
		arg := database.CreateGlovoRestaurantParams{
			Name:              rest.Name,
			Address:           rest.Address,
			DeliveryFee:       strconv.FormatFloat(rest.DeliveryFee, 'g', -1, 64),
			PhoneNumber:       sql.NullString{String: rest.PhoneNumber, Valid: true},
			GlovoApiStoreID:   int32(rest.GlovoApiStoreID),
			GlovoApiAddressID: int32(rest.GlovoApiAddressID),
			CreatedAt:         time.Now().UTC(),
			UpdatedAt:         time.Now().UTC(),
		}

		res, err := cfg.DB.CreateGlovoRestaurant(ctx, arg)
		if err != nil {
			log.Println("Couldn't post restaurant to DB :v", err)
			return err
		}
		if newRows, err := res.RowsAffected(); err != nil {
			log.Println("Couldn't access SQL result", err)
			rowsAffected += newRows
		}
		log.Printf("Created restaurant %v\n", rest.Name)
	}

	log.Printf("Created %v restaurants\n", rowsAffected)
	return nil
}

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

func FetchGlovoDishes(rest models.GlovoRestaurant, dishesURL string) ([]models.GlovoDish, error) {
	targetURL := strings.Replace(dishesURL, "{glovo_store_id}", strconv.Itoa(rest.GlovoApiStoreID), 1)
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
