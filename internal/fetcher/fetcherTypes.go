package fetcher

type glovoFiltersResponse struct {
	TopFilters []struct {
		FilterID     int    `json:"filterId"`
		CategoryID   int    `json:"categoryId"`
		CategoryName string `json:"categoryName"`
		CategorySlug string `json:"categorySlug"`
		FilterName   string `json:"filterName"`
		FilterSlug   string `json:"filterSlug"`
		Ranking      int    `json:"ranking"`
	} `json:"topFilters"`
}

type glovoRestaurantsResponse struct {
	Elements []struct {
		SingleData struct {
			StoreData struct {
				Store struct {
					Type            string  `json:"type"`
					ID              int     `json:"id"`
					Name            string  `json:"name"`
					Slug            string  `json:"slug"`
					Open            bool    `json:"open"`
					ServiceFee      float64 `json:"serviceFee"`
					DeliveryFeeInfo struct {
						Fee   float64 `json:"fee"`
						Style string  `json:"style"`
					} `json:"deliveryFeeInfo"`
					CategoryID  int         `json:"categoryId"`
					Distance    string      `json:"distance"`
					AddressID   int         `json:"addressId"`
					Location    interface{} `json:"location"`
					ItemsType   string      `json:"itemsType"`
					PhoneNumber string      `json:"phoneNumber"`
					Address     string      `json:"address"`
				} `json:"store"`
				Filters []struct {
					Name string `json:"name"`
				} `json:"filters"`
			} `json:"storeData"`
		} `json:"singleData"`
	} `json:"elements"`
}

type GlovoDishesResponse struct {
	Data struct {
		Body []struct {
			Data struct {
				Title    string `json:"title"`
				Slug     string `json:"slug"`
				Elements []struct {
					Type string `json:"type"`
					Data struct {
						ID             int64       `json:"id"`
						ExternalID     string      `json:"externalId"`
						StoreProductID interface{} `json:"storeProductId"`
						Name           string      `json:"name"`
						Description    string      `json:"description"`
						Price          float64     `json:"price"`
						PriceInfo      struct {
							Amount       float64 `json:"amount"`
							CurrencyCode string  `json:"currencyCode"`
							DisplayText  string  `json:"displayText"`
						} `json:"priceInfo"`
						AttributeGroups []struct {
							ID         int64  `json:"id"`
							ExternalID string `json:"externalId"`
							Name       string `json:"name"`
							Min        int    `json:"min"`
							Max        int    `json:"max"`
							Attributes []struct {
								ID          int64   `json:"id"`
								ExternalID  string  `json:"externalId"`
								PriceImpact float64 `json:"priceImpact"`
								PriceInfo   struct {
									Amount       float64 `json:"amount"`
									CurrencyCode string  `json:"currencyCode"`
									DisplayText  string  `json:"displayText"`
								} `json:"priceInfo"`
								Selected bool   `json:"selected"`
								Name     string `json:"name"`
							} `json:"attributes"`
							Position           int  `json:"position"`
							MultipleSelection  bool `json:"multipleSelection"`
							CollapsedByDefault bool `json:"collapsedByDefault"`
						} `json:"attributeGroups"`
						Promotion struct {
							ProductID   int64   `json:"productId"`
							PromotionID int     `json:"promotionId"`
							Title       string  `json:"title"`
							Type        string  `json:"type"`
							Percentage  float64 `json:"percentage"`
							Price       float64 `json:"price"`
							PriceInfo   struct {
								Amount       float64 `json:"amount"`
								CurrencyCode string  `json:"currencyCode"`
								DisplayText  string  `json:"displayText"`
							} `json:"priceInfo"`
							IsPrime bool   `json:"isPrime"`
							PromoID string `json:"promoId"`
						} `json:"promotion"`
						Promotions []struct {
							ProductID   int64   `json:"productId"`
							PromotionID int     `json:"promotionId"`
							Title       string  `json:"title"`
							Type        string  `json:"type"`
							Percentage  float64 `json:"percentage"`
							Price       float64 `json:"price"`
							PriceInfo   struct {
								Amount       float64 `json:"amount"`
								CurrencyCode string  `json:"currencyCode"`
								DisplayText  string  `json:"displayText"`
							} `json:"priceInfo"`
							IsPrime     bool   `json:"isPrime"`
							PromoID     string `json:"promoId"`
							Description string `json:"description,omitempty"`
						} `json:"promotions"`
						Indicators      []interface{} `json:"indicators"`
						Sponsored       bool          `json:"sponsored"`
						Restricted      bool          `json:"restricted"`
						ShowQuantifiers bool          `json:"showQuantifiers"`
					} `json:"data"`
				} `json:"elements"`
			} `json:"data"`
		} `json:"body"`
	} `json:"data"`
}
