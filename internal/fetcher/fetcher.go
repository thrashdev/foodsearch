package fetcher

import (
	"net/http"
)

type GlovoResponse struct {
	Title         string      `json:"title"`
	Version       int         `json:"version"`
	TotalElements int         `json:"totalElements"`
	Offset        int         `json:"offset"`
	RankingType   interface{} `json:"rankingType"`
	Elements      []struct {
		Type       string `json:"type"`
		SingleData struct {
			Type      string `json:"type"`
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
					Address                  string      `json:"address"`
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

func fetchGlovoItems(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Glovo-API-Version", "14")
	req.Header.Add("Glovo-App-Platform", "web")
	req.Header.Add("Glovo-App-Type", "customer")
	req.Header.Add("Glovo-App-Version", "7")
	req.Header.Add("Glovo-Location-City-Code", "BSK")

	return nil
}
