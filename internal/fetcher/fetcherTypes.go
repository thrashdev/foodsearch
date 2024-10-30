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

type glovoStoresResponse struct {
	Elements []struct {
		Type       string `json:"-"`
		SingleData struct {
			Type      string `json:"-"`
			StoreData struct {
				Store struct {
					Type        string      `json:"type"`
					ID          int         `json:"id"`
					Urn         string      `json:"urn"`
					Name        string      `json:"name"`
					Slug        string      `json:"slug"`
					FiscalName  interface{} `json:"fiscalName"`
					ImageID     string      `json:"imageId"`
					Open        bool        `json:"open"`
					EmulateOpen bool        `json:"emulateOpen"`
					McdPartner  bool        `json:"mcdPartner"`
					Food        bool        `json:"food"`
					CityCode    string      `json:"cityCode"`
					Scheduling  struct {
						Enabled bool        `json:"enabled"`
						Message interface{} `json:"message"`
					} `json:"scheduling"`
					ClosedStatusMessage interface{} `json:"closedStatusMessage"`
					NextOpeningTime     interface{} `json:"nextOpeningTime"`
					ServiceFee          float64     `json:"serviceFee"`
					DeliveryFeeInfo     struct {
						Fee   float64 `json:"fee"`
						Style string  `json:"style"`
					} `json:"deliveryFeeInfo"`
					CategoryID                  int           `json:"categoryId"`
					CartUniqueElements          interface{}   `json:"cartUniqueElements"`
					CartTotalElements           interface{}   `json:"cartTotalElements"`
					Note                        interface{}   `json:"note"`
					Distance                    string        `json:"distance"`
					AddressID                   int           `json:"addressId"`
					Location                    interface{}   `json:"location"`
					CustomDescriptionAllowed    bool          `json:"customDescriptionAllowed"`
					ProductsInformationText     interface{}   `json:"productsInformationText"`
					ProductsInformationLink     interface{}   `json:"productsInformationLink"`
					DeliveryNotAvailable        bool          `json:"deliveryNotAvailable"`
					DeliveryNotAvailableMessage interface{}   `json:"deliveryNotAvailableMessage"`
					SpecialRequirementsAllowed  bool          `json:"specialRequirementsAllowed"`
					EtaEnabled                  bool          `json:"etaEnabled"`
					AllergiesInformationAllowed bool          `json:"allergiesInformationAllowed"`
					LegalCheckboxRequired       bool          `json:"legalCheckboxRequired"`
					DataSharingRequested        bool          `json:"dataSharingRequested"`
					Marketplace                 bool          `json:"marketplace"`
					CashSupported               bool          `json:"cashSupported"`
					Promotions                  []interface{} `json:"promotions"`
					PrimeAvailable              bool          `json:"primeAvailable"`
					CanDisplayPrimeTierUpsell   bool          `json:"canDisplayPrimeTierUpsell"`
					CutleryRequestAllowed       bool          `json:"cutleryRequestAllowed"`
					RatingInfo                  struct {
						CardLabel        string      `json:"cardLabel"`
						DetailsLabel     interface{} `json:"detailsLabel"`
						TotalRatingLabel interface{} `json:"totalRatingLabel"`
						Icon             struct {
							LightImageID string `json:"lightImageId"`
							DarkImageID  string `json:"darkImageId"`
						} `json:"icon"`
						Color           interface{} `json:"color"`
						BackgroundColor struct {
							LightColorHex string `json:"lightColorHex"`
							DarkColorHex  string `json:"darkColorHex"`
						} `json:"backgroundColor"`
					} `json:"ratingInfo"`
					SelectedStrategyType interface{} `json:"selectedStrategyType"`
					SupportedStrategies  []struct {
						Type string `json:"type"`
					} `json:"supportedStrategies"`
					ItemsType                string      `json:"itemsType"`
					SuggestionKeywords       []string    `json:"suggestionKeywords"`
					PhoneNumber              string      `json:"phoneNumber"`
					Setress                  string      `json:"address"`
					ViewType                 string      `json:"viewType"`
					Sponsored                bool        `json:"sponsored"`
					FeesPricingCalculationID string      `json:"feesPricingCalculationId"`
					RankingScore             interface{} `json:"rankingScore"`
					Favorite                 bool        `json:"favorite"`
					EdenredEnabled           bool        `json:"edenredEnabled"`
				} `json:"store"`
				Filters []struct {
					Name string `json:"name"`
				} `json:"filters"`
			} `json:"storeData"`
			StoreProductsData            interface{} `json:"storeProductsData"`
			StoreWithProductCarouselData interface{} `json:"storeWithProductCarouselData"`
			BannerData                   interface{} `json:"bannerData"`
		} `json:"singleData"`
		GroupData interface{} `json:"groupData"`
	} `json:"elements"`
	RankingRequestID     interface{}   `json:"rankingRequestId"`
	ExactMatch           bool          `json:"exactMatch"`
	SearchSubVertical    string        `json:"searchSubVertical"`
	Actions              []interface{} `json:"actions"`
	WallAvailabilityData interface{}   `json:"wallAvailabilityData"`
}

type GlovoItemsResponse struct {
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
