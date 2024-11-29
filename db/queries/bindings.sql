-- name: BatchCreateRestaurantBinding :copyfrom
INSERT INTO restaurant_binding(id,glovo_restaurant_id, yandex_restaurant_id)
VALUES ($1, $2, $3);

-- name: GetAllRestaurantBindings :many
select * from restaurant_binding;
