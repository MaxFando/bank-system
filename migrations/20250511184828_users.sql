-- +goose Up
-- +goose StatementBegin
CREATE TABLE main.users
(
    id            SERIAL PRIMARY KEY,           -- Идентификатор пользователя
    email         VARCHAR(255) UNIQUE NOT NULL, -- Электронная почта (уникальна)
    password_hash VARCHAR(255)        NOT NULL, -- Хэш пароля
    first_name    VARCHAR(100),                 -- Имя
    last_name     VARCHAR(100),                 -- Фамилия
    date_of_birth DATE,                         -- Дата рождения
    created_at    TIMESTAMP DEFAULT NOW(),      -- Дата регистрации
    updated_at    TIMESTAMP DEFAULT NOW()       -- Дата последнего обновления
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS main.users;
-- +goose StatementEnd
