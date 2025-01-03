-- name: BatchCreateRestaurantBinding :copyfrom
INSERT INTO restaurant_binding(id,glovo_restaurant_id, yandex_restaurant_id)
VALUES ($1, $2, $3);

-- name: GetAllRestaurantBindings :many
select * from restaurant_binding;

-- name: GetRestaurantBindingsToUpdate :many
select distinct rb.* from restaurant_binding rb
left join dish_binding db on rb.id = db.restaurant_binding_id
where db.id is NULL;

-- name: BatchCreateDishBindings :copyfrom
INSERT INTO dish_binding(id, restaurant_binding_id, glovo_dish_id, yandex_dish_id)
VALUES ($1, $2, $3, $4);

-- name: GetAllRestaurants :many
select rb.ID, g.*, y.* from restaurant_binding rb
left join glovo_restaurant g on rb.glovo_restaurant_id = g.id
left join yandex_restaurant y on rb.yandex_restaurant_id = y.id;
