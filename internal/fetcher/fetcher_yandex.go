package fetcher

import (
	"fmt"
	"io"
	"net/http"
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

const search_slug = "search_restaurant"
const search_type = "quicksearch"

func fetchRestaurants(url string, loc Location, filter string) (payload []byte, err error) {
	query := YandexSearchQuery{Text: filter,
		Filters:  []YandexSearchFilter{YandexSearchFilter{Type: search_type, Slug: search_slug}},
		Selector: "all",
		Location: loc,
	}

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return []byte{}, fmt.Errorf("Couldn't create request, err: %v", err)
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, fmt.Errorf("Couldn't post request, err: %v", err)
	}
	defer resp.Body.Close()
	fmt.Println("Response Code: ", resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Couldn't read response body, err: %v", err)
	}

	return respBody, nil
}
