package models

import (
	"time"

	"github.com/google/uuid"
)

// not a part of DB
type Restaurant struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"store_name"`
	Address     string    `json:"address"`
	DeliveryFee float64   `json:"delivery_fee"`
	PhoneNumber string    `json:"phoneNumber"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// not a part of DB
type Dish struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	MaxDiscount float64   `json:"max_discount"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
	GlovoAPIDishID int `json:"id"`
}

type RestaurantBinding struct {
	GlovoRestaurantID  int
	YandexRestaurantID int
	GISRestaurantID    int
}

type DishBinding struct {
	GlovoDishID  int
	YandexDishID int
	GISDishID    int
}
