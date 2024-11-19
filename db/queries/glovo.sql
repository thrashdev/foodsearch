-- name: BatchCreateGlovoRestaurants :copyfrom
INSERT INTO glovo_restaurant(id, name, address, delivery_fee, phone_number, glovo_api_store_id, glovo_api_address_id, glovo_api_slug, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: BatchCreateGlovoDishes :copyfrom
INSERT INTO glovo_dish(id, name, description, price, discount, glovo_api_dish_id, glovo_restaurant_id, created_at, updated_at)
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

-- CREATE TABLE glovo_dish(
-- 	id UUID PRIMARY KEY,
-- 	name TEXT NOT NULL,
-- 	description TEXT NOT NULL,
-- 	price DECIMAL NOT NULL,
-- 	discount DECIMAL NOT NULL,
-- 	glovo_api_dish_id INTEGER NOT NULL,
-- 	glovo_restaurant_id UUID,
-- 	created_at TIMESTAMP NOT NULL,
-- 	updated_at TIMESTAMP NOT NULL,
-- 	CONSTRAINT fk_glovo_dish_glovo_restaurant_id
-- 	FOREIGN KEY (glovo_restaurant_id)
-- 	REFERENCES glovo_restaurant(id)
-- );
-- CREATE TABLE glovo_restaurant(
-- 	id UUID PRIMARY KEY,
-- 	name TEXT NOT NULL,
-- 	address TEXT NOT NULL,
-- 	delivery_fee DECIMAL NOT NULL,
-- 	phone_number TEXT,
-- 	glovo_api_store_id INTEGER NOT NULL,
-- 	glovo_api_address_id INTEGER NOT NULL,
-- 	created_at TIMESTAMP NOT NULL,
-- 	updated_at TIMESTAMP NOT NULL
-- );
--
-- CREATE TABLE restaurant_binding(
-- 	id UUID PRIMARY KEY,
-- 	glovo_restaurant_id UUID,
-- 	CONSTRAINT fk_restaurant_binding_glovo_restaurant_id
-- 	FOREIGN KEY (glovo_restaurant_id)
-- 	REFERENCES glovo_restaurant(id)
-- );
