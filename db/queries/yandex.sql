-- name: BatchCreateYandexRestaurants :copyfrom
INSERT INTO yandex_restaurant(id, name, address, delivery_fee, phone_number, yandex_api_slug,created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetYandexFilters :many
SELECT name FROM yandex_filters;

-- name: GetYandexRestaurantSlugs :many
SELECT yandex_api_slug FROM yandex_restaurant;

-- name: GetYandexRestaurant :one
select * from yandex_restaurant limit 1;
