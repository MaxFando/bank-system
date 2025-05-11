-- +goose Up
-- +goose StatementBegin
CREATE TABLE main.accounts
(
    id             SERIAL PRIMARY KEY,            -- Идентификатор счета
    user_id        INTEGER REFERENCES users (id), -- Внешний ключ на пользователя
    account_number VARCHAR(20) UNIQUE NOT NULL,   -- Номер счета
    balance        DECIMAL(15, 2) DEFAULT 0.00,   -- Баланс счета
    currency       VARCHAR(3)     DEFAULT 'RUB',  -- Валюта счета (фиксировано на RUB)
    account_type   VARCHAR(50),                   -- Тип счета (основной, сберегательный)
    created_at     TIMESTAMP      DEFAULT NOW(),  -- Дата создания счета
    updated_at     TIMESTAMP      DEFAULT NOW()   -- Дата последнего обновления
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS main.accounts;
-- +goose StatementEnd
