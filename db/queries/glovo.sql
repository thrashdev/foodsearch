-- name: BatchCreateGlovoRestaurants :copyfrom
INSERT INTO glovo_restaurant(id, name, address, delivery_fee, phone_number, glovo_api_store_id, glovo_api_address_id, glovo_api_slug, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: BatchCreateGlovoDishes :copyfrom
INSERT INTO glovo_dish(id, name, description, price, discounted_price, glovo_api_dish_id, glovo_restaurant_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: GetGlovoRestaurantsToUpdate :many
SELECT * FROM glovo_restaurant ORDER BY updated_at DESC
LIMIT $1;

-- name: GetAllGlovoRestaurants :many
select * from glovo_restaurant;

-- name: GetGlovoRestaurantNames :many
SELECT name FROM glovo_restaurant;

-- name: GetGlovoRestaurantsByName :many
SELECT * FROM glovo_restaurant
WHERE name = ANY(@restaurant_names::TEXT[]);

-- name: GetGlovoDishNames :many
select name from glovo_dish;

-- name: GetGlovoDishAPI_ID :many
select glovo_api_dish_id from glovo_dish;

-- name: GetAllGlovoDishes :many
select * from glovo_dish;

-- name: GetGlovoDishesForRestaurant :many
select * from glovo_dish
where glovo_restaurant_id = $1;

-- name: GetGlovoWeightedDishesForRestaurant :many
select * from glovo_dish
where name ~ '[0-9]+гр'
and glovo_restaurant_id = $1;
