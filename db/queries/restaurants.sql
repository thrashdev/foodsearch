-- name: CreateGlovoRestaurant :one
INSERT INTO glovo_restaurant(name, address, delivery_fee, phone_number, glovo_api_store_id, glovo_api_address_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	returning *;


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
