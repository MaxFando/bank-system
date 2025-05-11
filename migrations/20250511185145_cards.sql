-- +goose Up
-- +goose StatementBegin
CREATE TABLE main.cards
(
    id             SERIAL PRIMARY KEY,                    -- Идентификатор карты
    account_id     INTEGER REFERENCES main.accounts (id), -- Внешний ключ на счет
    encrypted_data TEXT NOT NULL,                         -- Зашифрованные данные карты
    hmac           TEXT NOT NULL,                         -- HMAC для проверки целостности данных
    created_at     TIMESTAMP DEFAULT NOW(),               -- Дата создания карты
    updated_at     TIMESTAMP DEFAULT NOW()                -- Дата последнего обновления
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS main.cards;
-- +goose StatementEnd
