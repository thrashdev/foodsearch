-- name: BatchCreateRestaurantBinding :copyfrom
INSERT INTO restaurant_binding(id,glovo_restaurant_id, yandex_restaurant_id)
VALUES ($1, $2, $3);

-- name: GetAllRestaurantBindings :many
select * from restaurant_binding;

-- name: BatchCreateDishBindings :copyfrom
INSERT INTO dish_binding(id, restaurant_binding_id, glovo_dish_id, yandex_dish_id)
VALUES ($1, $2, $3, $4);
