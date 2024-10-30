package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/thrashdev/foodsearch/internal/models"
)

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

func fetchGlovoStoresByFilter(baseURL string, filter string) (restaurants []string, err error) {
	fullURL := baseURL + "&filter=" + filter
	respBody, err := fetchByUrl(fullURL)

	var glovoResp glovoStoresResponse
	err = json.Unmarshal(respBody, &glovoResp)
	if err != nil {
		log.Println(fullURL)
		return []string{}, err
	}

	for _, item := range glovoResp.Elements {
		restaurants = append(restaurants, item.SingleData.StoreData.Store.Name)
	}

	return restaurants, nil

}

func FetchGlovoItems(baseURL string, filtersURL string) (allRestaurants []string, err error) {
	filters, err := fetchGlovoFilters(filtersURL)
	if err != nil {
		return []string{}, fmt.Errorf("Couldn't get filters, err: %v", err)
	}

	for _, f := range filters {
		restaurantsByFilter, err := fetchGlovoStoresByFilter(baseURL, url.QueryEscape(f))
		if err != nil {
			return []string{}, fmt.Errorf("Couldn't fetch by filter: %s. Error :%v", f, err)
		}

		allRestaurants = append(allRestaurants, restaurantsByFilter...)
	}
	return allRestaurants, nil
}

func GetGlovoStores(resp glovoStoresResponse) ([]models.GlovoStore, error) {
	result := []models.GlovoStore{}
	for _, item := range resp.Elements {
		store := models.GlovoStore{ID: uuid.New(),
			GlovoStoreID:    item.SingleData.StoreData.Store.ID,
			GlovoAddressID:  item.SingleData.StoreData.Store.AddressID,
			Name:            item.SingleData.StoreData.Store.Name,
			DeliveryFee:     item.SingleData.StoreData.Store.ServiceFee,
			DeliveryFeeInfo: item.SingleData.StoreData.Store.DeliveryFeeInfo,
		}

		result = append(result, store)
	}

	return result, nil
}
