package models

import (
	"github.com/google/uuid"
)

type GlovoStore struct {
	ID              uuid.UUID `json:"id"`
	GlovoStoreID    int       `json:"glovo_id"`
	GlovoAddressID  int       `json:"glovo_address_id"`
	Name            string    `json:"store_name"`
	DeliveryFee     float64   `json:"serviceFee"`
	DeliveryFeeInfo struct {
		Fee   float64 `json:"fee"`
		Style string  `json:"style"`
	} `json:"deliveryFeeInfo"`
}

type GlovoDish struct {
	GlovoID     int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	PriceInfo   struct {
		Amount       float64 `json:"amount"`
		CurrencyCode string  `json:"currencyCode"`
		DisplayText  string  `json:"displayText"`
	} `json:"priceInfo"`
	CuisineID int `json:"cuisine_id"`
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
	} `json:"promotion"`
}
