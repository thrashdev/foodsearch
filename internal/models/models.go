package models

import (
	"time"

	"github.com/google/uuid"
)

// not a part of DB
type Restaurant struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"store_name"`
	Address     *string   `json:"address"`
	DeliveryFee *float64  `json:"delivery_fee"`
	PhoneNumber *string   `json:"phoneNumber"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// not a part of DB
type Dish struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Price           float64   `json:"price"`
	DiscountedPrice float64   `json:"max_discount"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type DishCategory struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type GlovoRestaurant struct {
	Restaurant
	GlovoApiStoreID   int    `json:"glovo_id"`
	GlovoApiAddressID int    `json:"glovo_address_id"`
	GlovoApiSlug      string `json:"glovo_api_slug"`
}

type GlovoDish struct {
	Dish
	GlovoAPIDishID    int       `json:"id"`
	GlovoRestaurantID uuid.UUID `json:"glovo_restaurant_id"`
}

type YandexRestaurant struct {
	Restaurant
	YandexApiSlug string
}

type YandexDish struct {
	Dish
	YandexRestaurantID uuid.UUID `json:"yandex_restaurant_id"`
	YandexApiID        int
}

type RestaurantBinding struct {
	GlovoRestaurantID  uuid.UUID
	YandexRestaurantID uuid.UUID
}

type DishBinding struct {
	GlovoDishID  uuid.UUID
	YandexDishID uuid.UUID
}
