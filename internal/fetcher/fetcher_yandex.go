package fetcher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/thrashdev/foodsearch/internal/config"
	"github.com/thrashdev/foodsearch/internal/database"
	"github.com/thrashdev/foodsearch/internal/models"
	"github.com/thrashdev/foodsearch/internal/utils"
)

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type YandexSearchQuery struct {
	Text     string               `json:"text"`
	Filters  []YandexSearchFilter `json:"filters"`
	Selector string               `json:"selector"`
	Location Location             `json:"location"`
}

type YandexSearchFilter struct {
	Type string `json:"type"`
	Slug string `json:"slug"`
}

var yandex_categories_blacklist = []string{"напитки", "хлеб", "закуски", "дополнительно", "соусы",
	"горячие напитки", "холодные напитки", "соки", "бутылочное", "к пиву", "пиво бутылочное", "сигареты",
	"лимонады и айс-ти", "кофе", "алкогольные напитки", "вино", "добавки", "алкоголь", "пиво",
	"гарниры", "гарнир", "чаи", "полуфабрикаты", "безалкогольные напитки"}

const search_slug = "search_restaurant"
const search_type = "quickfilter"

const search_url_slug_token = "{restaurant_slug}"
const search_url_longitude_token = "{longitude}"
const search_url_latitude_token = "{latitude}"

func InitYandex(cfg *config.ServiceConfig) {
	startupCommands := []startupCommand{
		CreateNewYandexRestaurants,
		CreateNewYandexDishes,
	}

	for _, cmd := range startupCommands {
		res, err := cmd(cfg)
		if err != nil {
			cfg.Logger.DPanic("Error during yandex startup", "error", err)
		}
		res.print()
	}
}

func restaurantDifferenceYandex(restaurants []models.YandexRestaurant, slugs []string) []models.YandexRestaurant {
	mb := make(map[string]struct{}, len(slugs))
	for _, slug := range slugs {
		mb[slug] = struct{}{}
	}
	var diff []models.YandexRestaurant
	for _, rest := range restaurants {
		if _, found := mb[rest.YandexApiSlug]; !found {
			diff = append(diff, rest)
		}
	}
	return diff

}

func removeDuplicateYandexRestaurants(rests []models.YandexRestaurant) []models.YandexRestaurant {
	extracted := make(map[string]struct{})
	result := []models.YandexRestaurant{}
	for _, r := range rests {
		_, ok := extracted[r.Name]
		if ok {
			continue
		}
		result = append(result, r)
		extracted[r.Name] = struct{}{}
	}
	return result
}

func removeDuplicateYandexDishes(dishes []models.YandexDish) []models.YandexDish {
	deduped := make(map[int]struct{})
	result := []models.YandexDish{}
	for _, d := range dishes {
		_, exists := deduped[d.YandexApiID]
		if exists {
			continue
		}
		result = append(result, d)
		deduped[d.YandexApiID] = struct{}{}
	}
	return result
}

func dishYandexApiIdDifference(dishes []models.YandexDish, ids []int32) []models.YandexDish {
	mb := make(map[int32]struct{}, len(ids))
	for _, id := range ids {
		mb[id] = struct{}{}
	}
	uniqueDishes := []models.YandexDish{}
	for _, dish := range dishes {
		if _, found := mb[int32(dish.YandexApiID)]; !found {
			uniqueDishes = append(uniqueDishes, dish)
		}
	}
	return uniqueDishes
}

func CreateNewYandexRestaurants(cfg *config.ServiceConfig) (DBActionResult, error) {
	ctx := context.Background()
	filters, err := cfg.DB.GetYandexFilters(ctx)
	if err != nil {
		return DBActionResult{}, fmt.Errorf("Couldn't fetch yandex filters, err: %w", err)
	}

	restsWithDuplicates := []models.YandexRestaurant{}
	for _, f := range filters {
		rests, err := FetchYandexRestaurants(cfg, f)
		if err != nil {
			cfg.Logger.Warnf("Error encountered while fetching yandex restaurants by filter: %v, err: %v", f, err)
			continue
		}
		restsWithDuplicates = append(restsWithDuplicates, rests...)
	}

	slugs, err := cfg.DB.GetYandexRestaurantSlugs(ctx)
	if err != nil {
		return DBActionResult{}, fmt.Errorf("Couldn't fetch yandex restaurant slugs, err: %w", err)
	}
	rests := removeDuplicateYandexRestaurants(restsWithDuplicates)
	newRests := restaurantDifferenceYandex(rests, slugs)
	rowsAffected, err := createYandexRestaurants(cfg, newRests)
	if err != nil {
		return DBActionResult{}, fmt.Errorf("Couldn't post yandex restaurants to DB: %w", err)
	}
	result := DBActionResult{}
	result.records = append(result.records, makeDBActionResultRecord("Created %v new yandex restaurants", rowsAffected))
	return result, nil

}

func createYandexRestaurants(cfg *config.ServiceConfig, rests []models.YandexRestaurant) (rowsAffected int64, err error) {
	args := []database.BatchCreateYandexRestaurantsParams{}
	for _, r := range rests {
		arg := utils.YandexRestModelToDB(r)
		args = append(args, arg)
	}

	ctx := context.Background()
	rowsAffected, err = cfg.DB.BatchCreateYandexRestaurants(ctx, args)
	if err != nil {
		return 0, fmt.Errorf("Error creating Yandex restaurants in DB: %v", err)
	}
	return rowsAffected, err

}

func CreateNewYandexDishes(cfg *config.ServiceConfig) (DBActionResult, error) {
	ctx := context.Background()
	rests, err := cfg.DB.GetAllYandexRestaurants(ctx)
	if err != nil {
		return DBActionResult{}, fmt.Errorf("Error fetching yandex restaurants from DB: %w", err)
	}

	dishesResp := []models.YandexDish{}
	for _, dbRest := range rests {
		rest := utils.YandexRestDBtoModel(dbRest)
		dishesPerRest := FetchYandexDishes(cfg, rest)
		dishesResp = append(dishesResp, dishesPerRest...)
	}
	dedupedDishes := removeDuplicateYandexDishes(dishesResp)
	yandexApiIDS, err := cfg.DB.GetYandexDishApiIDS(ctx)
	if err != nil {
		return DBActionResult{}, fmt.Errorf("Couldn't fetch yandex api ids from DB")
	}
	dishes := dishYandexApiIdDifference(dedupedDishes, yandexApiIDS)

	rowsAffected, err := postYandexDishesToDB(cfg, dishes)
	if err != nil {
		return DBActionResult{}, fmt.Errorf("Couldn't post new yandex dishes to DB, %v", err)
	}
	result := DBActionResult{}
	result.records = append(result.records, makeDBActionResultRecord("Created %v new dishes", rowsAffected))
	return result, nil

}

func createNewYandexDishes(cfg *config.ServiceConfig, restID uuid.UUID) {

}

func postYandexDishesToDB(cfg *config.ServiceConfig, dishes []models.YandexDish) (rowsAffected int64, err error) {
	args := []database.BatchCreateYandexDishesParams{}
	for _, d := range dishes {
		arg := utils.YandexDishModelToDB(d)
		args = append(args, arg)
	}
	ctx := context.Background()
	rowsAffected, err = cfg.DB.BatchCreateYandexDishes(ctx, args)
	if err != nil {
		return 0, fmt.Errorf("Couldn't create yandex dishes in DB: %v", err)
	}
	return rowsAffected, nil
}

func FetchYandexRestaurants(cfg *config.ServiceConfig, filter string) (allRestaurants []models.YandexRestaurant, err error) {
	query := YandexSearchQuery{Text: filter,
		Filters:  []YandexSearchFilter{YandexSearchFilter{Type: search_type, Slug: search_slug}},
		Selector: "all",
		Location: Location(cfg.Yandex.Loc),
	}

	buf, err := json.Marshal(query)
	if err != nil {
		fmt.Println("error while marhsalling json: %v", err)
		return nil, err
	}

	body := bytes.NewBuffer(buf)

	req, err := http.NewRequest("POST", cfg.Yandex.SearchURL, body)
	if err != nil {
		return nil, fmt.Errorf("Couldn't create request, err: %v", err)
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Couldn't post request, err: %v", err)
	}
	defer resp.Body.Close()
	respB := []byte{}
	if resp.StatusCode > 200 {
		respB, err = io.ReadAll(resp.Body)
	}
	fmt.Printf("Fetched Yandex Restaurant | Response Code: %v, Body: %v\n", resp.StatusCode, string(respB))

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Couldn't read response body, err: %v", err)
	}

	var response YandexSearchResponse
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return nil, fmt.Errorf("Error during JSON unmarshalling: %v", err)
	}

	for _, block := range response.Blocks {
		if block.Type != "places" {
			continue
		}

		for _, rest := range block.Payload {
			allRestaurants = append(allRestaurants, models.YandexRestaurant{
				Restaurant: models.Restaurant{
					ID:          uuid.New(),
					Name:        rest.Title,
					Address:     nil,
					DeliveryFee: nil,
					PhoneNumber: nil,
					CreatedAt:   time.Now().UTC(),
					UpdatedAt:   time.Now().UTC(),
				},
				YandexApiSlug: rest.Slug,
			})
		}
	}

	return allRestaurants, nil
}

func FetchYandexDishes(cfg *config.ServiceConfig, rest models.YandexRestaurant) []models.YandexDish {
	url := strings.Replace(cfg.Yandex.RestaurantMenuURL, search_url_slug_token, rest.YandexApiSlug, 1)
	url = strings.Replace(url, search_url_latitude_token, strconv.FormatFloat(cfg.Yandex.Loc.Latitude, 'f', -1, 64), 1)
	url = strings.Replace(url, search_url_longitude_token, strconv.FormatFloat(cfg.Yandex.Loc.Longitude, 'f', -1, 64), 1)
	fmt.Println("URL: ", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Couldn't fetch restaurant menu: %v", err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Couldn't read response body, err: %v", err)
	}

	var yandexResp YandexRestaurantMenuResponse
	err = json.Unmarshal(respBody, &yandexResp)
	if err != nil {
		log.Fatalf("Error parsing YandexMenuResponse: %v", err)
	}

	dishes := []models.YandexDish{}
	for _, ct := range yandexResp.Payload.Categories {
		categoryName := strings.ToLower(ct.Name)
		if utils.SliceContains(yandex_categories_blacklist, categoryName) {
			continue
		}
		// fmt.Printf("Yandex Category: %v\n", categoryName)
		for _, item := range ct.Items {
			dish := models.YandexDish{
				Dish: models.Dish{
					ID:              uuid.New(),
					Name:            item.Name,
					Description:     item.Description,
					Price:           float64(item.Price),
					DiscountedPrice: float64(item.PromoPrice),
					CreatedAt:       time.Now().UTC(),
					UpdatedAt:       time.Now().UTC(),
				},
				YandexRestaurantID: rest.ID,
				YandexApiID:        int(item.ID),
			}
			// fmt.Printf("	- %v\n", item.Name)
			dishes = append(dishes, dish)
		}

	}

	return dishes

}
