// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: restaurants.sql

package database

import (
	"context"
	"database/sql"
	"time"
)

const createGlovoRestaurant = `-- name: CreateGlovoRestaurant :one
INSERT INTO glovo_restaurant(name, address, delivery_fee, phone_number, glovo_api_store_id, glovo_api_address_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	returning id, name, address, delivery_fee, phone_number, glovo_api_store_id, glovo_api_address_id, created_at, updated_at
`

type CreateGlovoRestaurantParams struct {
	Name              string
	Address           string
	DeliveryFee       string
	PhoneNumber       sql.NullString
	GlovoApiStoreID   int32
	GlovoApiAddressID int32
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (q *Queries) CreateGlovoRestaurant(ctx context.Context, arg CreateGlovoRestaurantParams) (GlovoRestaurant, error) {
	row := q.db.QueryRowContext(ctx, createGlovoRestaurant,
		arg.Name,
		arg.Address,
		arg.DeliveryFee,
		arg.PhoneNumber,
		arg.GlovoApiStoreID,
		arg.GlovoApiAddressID,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i GlovoRestaurant
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Address,
		&i.DeliveryFee,
		&i.PhoneNumber,
		&i.GlovoApiStoreID,
		&i.GlovoApiAddressID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}