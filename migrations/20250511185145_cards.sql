-- +goose Up
-- +goose StatementBegin
CREATE TABLE main.cards
(
    id              SERIAL PRIMARY KEY,                    -- Идентификатор карты
    account_id      INTEGER REFERENCES main.accounts (id), -- Внешний ключ на счет
    card_number     VARCHAR(16) UNIQUE NOT NULL,           -- Номер карты
    expiration_date DATE,                                  -- Срок действия карты
    cvv             VARCHAR(3),                            -- CVV код
    status          VARCHAR(50),                           -- Статус карты (активная, заблокирована)
    created_at      TIMESTAMP DEFAULT NOW(),               -- Дата создания карты
    updated_at      TIMESTAMP DEFAULT NOW()                -- Дата последнего обновления
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS main.cards;
-- +goose StatementEnd
