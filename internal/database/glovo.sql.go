// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: glovo.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type BatchCreateGlovoDishesParams struct {
	ID                pgtype.UUID
	Name              string
	Description       string
	Price             pgtype.Numeric
	DiscountedPrice   pgtype.Numeric
	GlovoApiDishID    int32
	GlovoRestaurantID pgtype.UUID
	CreatedAt         pgtype.Timestamp
	UpdatedAt         pgtype.Timestamp
}

type BatchCreateGlovoRestaurantsParams struct {
	ID                pgtype.UUID
	Name              string
	Address           string
	DeliveryFee       pgtype.Numeric
	PhoneNumber       pgtype.Text
	GlovoApiStoreID   int32
	GlovoApiAddressID int32
	GlovoApiSlug      string
	CreatedAt         pgtype.Timestamp
	UpdatedAt         pgtype.Timestamp
}

const getAllGlovoDishes = `-- name: GetAllGlovoDishes :many
select id, name, description, price, discounted_price, glovo_api_dish_id, glovo_restaurant_id, created_at, updated_at from glovo_dish
`

func (q *Queries) GetAllGlovoDishes(ctx context.Context) ([]GlovoDish, error) {
	rows, err := q.db.Query(ctx, getAllGlovoDishes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GlovoDish
	for rows.Next() {
		var i GlovoDish
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Price,
			&i.DiscountedPrice,
			&i.GlovoApiDishID,
			&i.GlovoRestaurantID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllGlovoRestaurants = `-- name: GetAllGlovoRestaurants :many
select id, name, address, delivery_fee, phone_number, glovo_api_store_id, glovo_api_address_id, glovo_api_slug, created_at, updated_at from glovo_restaurant
`

func (q *Queries) GetAllGlovoRestaurants(ctx context.Context) ([]GlovoRestaurant, error) {
	rows, err := q.db.Query(ctx, getAllGlovoRestaurants)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GlovoRestaurant
	for rows.Next() {
		var i GlovoRestaurant
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Address,
			&i.DeliveryFee,
			&i.PhoneNumber,
			&i.GlovoApiStoreID,
			&i.GlovoApiAddressID,
			&i.GlovoApiSlug,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getGlovoDishAPI_ID = `-- name: GetGlovoDishAPI_ID :many
select glovo_api_dish_id from glovo_dish
`

func (q *Queries) GetGlovoDishAPI_ID(ctx context.Context) ([]int32, error) {
	rows, err := q.db.Query(ctx, getGlovoDishAPI_ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int32
	for rows.Next() {
		var glovo_api_dish_id int32
		if err := rows.Scan(&glovo_api_dish_id); err != nil {
			return nil, err
		}
		items = append(items, glovo_api_dish_id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getGlovoDishNames = `-- name: GetGlovoDishNames :many
select name from glovo_dish
`

func (q *Queries) GetGlovoDishNames(ctx context.Context) ([]string, error) {
	rows, err := q.db.Query(ctx, getGlovoDishNames)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		items = append(items, name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getGlovoDishesForRestaurant = `-- name: GetGlovoDishesForRestaurant :many
select id, name, description, price, discounted_price, glovo_api_dish_id, glovo_restaurant_id, created_at, updated_at from glovo_dish
where glovo_restaurant_id = $1
`

func (q *Queries) GetGlovoDishesForRestaurant(ctx context.Context, glovoRestaurantID pgtype.UUID) ([]GlovoDish, error) {
	rows, err := q.db.Query(ctx, getGlovoDishesForRestaurant, glovoRestaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GlovoDish
	for rows.Next() {
		var i GlovoDish
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Price,
			&i.DiscountedPrice,
			&i.GlovoApiDishID,
			&i.GlovoRestaurantID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getGlovoRestaurantNames = `-- name: GetGlovoRestaurantNames :many
SELECT name FROM glovo_restaurant
`

func (q *Queries) GetGlovoRestaurantNames(ctx context.Context) ([]string, error) {
	rows, err := q.db.Query(ctx, getGlovoRestaurantNames)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		items = append(items, name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getGlovoRestaurantsByName = `-- name: GetGlovoRestaurantsByName :many
SELECT id, name, address, delivery_fee, phone_number, glovo_api_store_id, glovo_api_address_id, glovo_api_slug, created_at, updated_at FROM glovo_restaurant
WHERE name = ANY($1::TEXT[])
`

func (q *Queries) GetGlovoRestaurantsByName(ctx context.Context, restaurantNames []string) ([]GlovoRestaurant, error) {
	rows, err := q.db.Query(ctx, getGlovoRestaurantsByName, restaurantNames)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GlovoRestaurant
	for rows.Next() {
		var i GlovoRestaurant
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Address,
			&i.DeliveryFee,
			&i.PhoneNumber,
			&i.GlovoApiStoreID,
			&i.GlovoApiAddressID,
			&i.GlovoApiSlug,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getGlovoRestaurantsToUpdate = `-- name: GetGlovoRestaurantsToUpdate :many
SELECT id, name, address, delivery_fee, phone_number, glovo_api_store_id, glovo_api_address_id, glovo_api_slug, created_at, updated_at FROM glovo_restaurant ORDER BY updated_at DESC
LIMIT $1
`

func (q *Queries) GetGlovoRestaurantsToUpdate(ctx context.Context, limit int32) ([]GlovoRestaurant, error) {
	rows, err := q.db.Query(ctx, getGlovoRestaurantsToUpdate, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GlovoRestaurant
	for rows.Next() {
		var i GlovoRestaurant
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Address,
			&i.DeliveryFee,
			&i.PhoneNumber,
			&i.GlovoApiStoreID,
			&i.GlovoApiAddressID,
			&i.GlovoApiSlug,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getGlovoWeightedDishesForRestaurant = `-- name: GetGlovoWeightedDishesForRestaurant :many
select id, name, description, price, discounted_price, glovo_api_dish_id, glovo_restaurant_id, created_at, updated_at from glovo_dish
where name ~ '[0-9]+гр'
and glovo_restaurant_id = $1
`

func (q *Queries) GetGlovoWeightedDishesForRestaurant(ctx context.Context, glovoRestaurantID pgtype.UUID) ([]GlovoDish, error) {
	rows, err := q.db.Query(ctx, getGlovoWeightedDishesForRestaurant, glovoRestaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GlovoDish
	for rows.Next() {
		var i GlovoDish
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Price,
			&i.DiscountedPrice,
			&i.GlovoApiDishID,
			&i.GlovoRestaurantID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
