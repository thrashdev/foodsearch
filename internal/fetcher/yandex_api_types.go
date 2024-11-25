package fetcher

type YandexSearchResponse struct {
	Header struct {
		Text string `json:"text"`
	} `json:"header"`
	Blocks []struct {
		Type    string `json:"type"`
		Payload []struct {
			Slug  string `json:"slug"`
			Title string `json:"title"`
			Tags  []struct {
				Title string `json:"title"`
			} `json:"tags"`
		} `json:"payload"`
	} `json:"blocks"`
}

type YandexRestaurantMenuResponse struct {
	Payload struct {
		Categories []struct {
			ID        int64  `json:"id"`
			Name      string `json:"name"`
			Available bool   `json:"available"`
			Items     []struct {
				ID           int64  `json:"id"`
				Name         string `json:"name"`
				Description  string `json:"description"`
				Available    bool   `json:"available"`
				InStock      any    `json:"inStock"`
				Price        int    `json:"price"`
				DecimalPrice string `json:"decimalPrice"`
				PromoPrice   int    `json:"promoPrice"`
				PromoTypes   []any  `json:"promoTypes"`
			} `json:"items"`
		} `json:"categories"`
	} `json:"payload"`
}
