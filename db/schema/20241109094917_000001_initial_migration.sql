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
	discounted_price DECIMAL NOT NULL,
	glovo_api_dish_id INTEGER NOT NULL UNIQUE,
	glovo_restaurant_id UUID NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	CONSTRAINT fk_glovo_dish_glovo_restaurant_id
	FOREIGN KEY (glovo_restaurant_id)
	REFERENCES glovo_restaurant(id)
);

CREATE TABLE yandex_restaurant(
	id UUID PRIMARY KEY,
	name TEXT NOT NULL,
	address TEXT ,
	delivery_fee DECIMAL,
	phone_number TEXT,
	yandex_api_slug TEXT NOT NULL UNIQUE,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL
);
	-- ID          uuid.UUID `json:"id"`
	-- Name        string    `json:"store_name"`
	-- Address     *string   `json:"address"`
	-- DeliveryFee *float64  `json:"delivery_fee"`
	-- PhoneNumber *string   `json:"phoneNumber"`
	-- CreatedAt   time.Time `json:"created_at"`
	-- UpdatedAt   time.Time `json:"updated_at"`
	-- YandexApiSlug string

CREATE TABLE restaurant_binding(
	id UUID PRIMARY KEY,
	glovo_restaurant_id UUID,
	yandex_restaurant_id UUID,

	CONSTRAINT fk_restaurant_binding_glovo_restaurant_id
	FOREIGN KEY (glovo_restaurant_id)
	REFERENCES glovo_restaurant(id),

	CONSTRAINT fk_restaurant_binding_yandex_restaurant_id
	FOREIGN KEY (yandex_restaurant_id)
	REFERENCES yandex_restaurant(id)
);

CREATE TABLE dish_binding(
	id UUID PRIMARY KEY,
	glovo_dish_id UUID,

	CONSTRAINT fk_dish_binding_glovo_dish_id
	FOREIGN KEY (glovo_dish_id)
	REFERENCES glovo_dish(id)
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
alter table dish_binding
drop constraint fk_dish_binding_glovo_dish_id;

alter table restaurant_binding
drop constraint fk_restaurant_binding_yandex_restaurant_id;

alter table restaurant_binding
drop constraint fk_restaurant_binding_glovo_restaurant_id;

alter table glovo_dish
drop constraint fk_glovo_dish_glovo_restaurant_id;

drop table restaurant_binding;
drop table dish_binding;
drop table glovo_restaurant;
drop table glovo_dish;
drop table yandex_restaurant;
-- +goose StatementEnd
