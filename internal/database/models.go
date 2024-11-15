// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type DishBinding struct {
	ID          uuid.UUID
	GlovoDishID uuid.NullUUID
}

type GlovoDish struct {
	ID                uuid.UUID
	Name              string
	Description       string
	Price             string
	Discount          string
	GlovoApiDishID    int32
	GlovoRestaurantID uuid.NullUUID
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type GlovoRestaurant struct {
	ID                uuid.UUID
	Name              string
	Address           string
	DeliveryFee       string
	PhoneNumber       sql.NullString
	GlovoApiStoreID   int32
	GlovoApiAddressID int32
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type RestaurantBinding struct {
	ID                uuid.UUID
	GlovoRestaurantID uuid.NullUUID
}
