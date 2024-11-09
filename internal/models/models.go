package models

import (
	"github.com/google/uuid"
)

type Restaurant struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"store_name"`
	Address     string    `json:"address"`
	DeliveryFee float64   `json:"delivery_fee"`
	PhoneNumber string    `json:"phoneNumber"`
}

type Dish struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	MaxDiscount float64   `json:"max_discount"`
}

type DishCategory struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type GlovoRestaurant struct {
	Restaurant
	GlovoStoreID   int `json:"glovo_id"`
	GlovoAddressID int `json:"glovo_address_id"`
}

type GlovoDish struct {
	Dish
	GlovoID int `json:"id"`
}
