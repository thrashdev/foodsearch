package models

import (
	"github.com/google/uuid"
)

type Store struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"store_name"`
	Address     string    `json:"address"`
	DeliveryFee float64   `json:"delivery_fee"`
}

type GlovoStore struct {
	Store
	GlovoStoreID   int `json:"glovo_id"`
	GlovoAddressID int `json:"glovo_address_id"`
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
