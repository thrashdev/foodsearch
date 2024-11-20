-- +goose Up
-- +goose StatementBegin
CREATE TABLE glovo_restaurant(
	id UUID PRIMARY KEY,
	name TEXT NOT NULL,
	address TEXT NOT NULL,
	delivery_fee DECIMAL NOT NULL,
	phone_number TEXT,
	glovo_api_store_id INTEGER NOT NULL,
	glovo_api_address_id INTEGER NOT NULL,
	glovo_api_slug TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL
);

CREATE TABLE glovo_dish(
	id UUID PRIMARY KEY,
	name TEXT NOT NULL,
	description TEXT NOT NULL,
	price DECIMAL NOT NULL,
	discount DECIMAL NOT NULL,
	glovo_api_dish_id INTEGER NOT NULL UNIQUE,
	glovo_restaurant_id UUID NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	CONSTRAINT fk_glovo_dish_glovo_restaurant_id
	FOREIGN KEY (glovo_restaurant_id)
	REFERENCES glovo_restaurant(id)
);

CREATE TABLE dish_binding(
	id UUID PRIMARY KEY,
	glovo_dish_id UUID,
	CONSTRAINT fk_dish_binding_glovo_dish_id
	FOREIGN KEY (glovo_dish_id)
	REFERENCES glovo_dish(id)
);

CREATE TABLE restaurant_binding(
	id UUID PRIMARY KEY,
	glovo_restaurant_id UUID,
	CONSTRAINT fk_restaurant_binding_glovo_restaurant_id
	FOREIGN KEY (glovo_restaurant_id)
	REFERENCES glovo_restaurant(id)
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
alter table dish_binding
drop constraint fk_dish_binding_glovo_dish_id;
alter table restaurant_binding
drop constraint fk_restaurant_binding_glovo_restaurant_id;
alter table glovo_dish
drop constraint fk_glovo_dish_glovo_restaurant_id;
drop table restaurant_binding;
drop table dish_binding;
drop table glovo_restaurant;
drop table glovo_dish;
-- +goose StatementEnd
