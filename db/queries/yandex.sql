-- name: BatchCreateYandexRestaurants :copyfrom
INSERT INTO yandex_restaurant(id, name, address, delivery_fee, phone_number, yandex_api_slug,created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: BatchCreateYandexDishes :copyfrom
INSERT INTO yandex_dish(id, name, description, price, discounted_price, yandex_restaurant_id, yandex_api_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: GetYandexFilters :many
SELECT name FROM yandex_filters;

-- name: GetYandexRestaurantSlugs :many
SELECT yandex_api_slug FROM yandex_restaurant;

-- name: GetYandexRestaurant :one
select * from yandex_restaurant limit 1;

-- name: GetAllYandexRestaurants :many
select * from yandex_restaurant;

-- name: GetAllYandexDishes :many
select * from yandex_dish;

-- name: GetYandexDishApiIDS :many
select yandex_api_id from yandex_dish;

