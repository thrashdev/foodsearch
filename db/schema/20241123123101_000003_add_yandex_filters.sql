-- +goose Up
-- +goose StatementBegin
CREATE TABLE yandex_filters(
id UUID PRIMARY KEY,
name TEXT NOT NULL UNIQUE
);

INSERT INTO yandex_filters(id, name)
VALUES
(gen_random_uuid(), 'Шаурма'),
(gen_random_uuid(), 'Бургер'),
(gen_random_uuid(), 'Хачапури'),
(gen_random_uuid(), 'Лагман'),
(gen_random_uuid(), 'Манты'),
(gen_random_uuid(), 'Плов'),
(gen_random_uuid(), 'Курица'),
(gen_random_uuid(), 'Картофель фри'),
(gen_random_uuid(), 'Суши'),
(gen_random_uuid(), 'Торт'),
(gen_random_uuid(), 'Кекс'),
(gen_random_uuid(), 'Десерт'),
(gen_random_uuid(), 'Пицца'),
(gen_random_uuid(), 'Завтрак'),
(gen_random_uuid(), 'Стейк'),
(gen_random_uuid(), 'Шашлык'),
(gen_random_uuid(), 'Паста'),
(gen_random_uuid(), 'Кебаб');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE yandex_filters;
-- +goose StatementEnd
