-- +goose Up
-- +goose StatementBegin
CREATE INDEX glovo_dish_name_idx on glovo_dish(name);
CREATE INDEX glovo_restaurant_name_idx on glovo_restaurant(name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX glovo_dish_name_idx;
DROP INDEX glovo_restaurant_name_idx;
-- +goose StatementEnd
