-- name: CreateGlovoRestaurant :one
INSERT INTO glovo_restaurant(name, address, delivery_fee, phone_number, glovo_api_store_id, glovo_api_address_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	returning *;

-- name: CreateGlovoDish :one
INSERT INTO glovo_dish(name, description, price, discount, glovo_api_dish_id, glovo_restaurant_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	returning *;


CREATE TABLE glovo_dish(
	id UUID PRIMARY KEY,
	name TEXT NOT NULL,
	description TEXT NOT NULL,
	price DECIMAL NOT NULL,
	discount DECIMAL NOT NULL,
	glovo_api_dish_id INTEGER NOT NULL,
	glovo_restaurant_id UUID,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	CONSTRAINT fk_glovo_dish_glovo_restaurant_id
	FOREIGN KEY (glovo_restaurant_id)
	REFERENCES glovo_restaurant(id)
);
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
