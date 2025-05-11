-- +goose Up
-- +goose StatementBegin
CREATE TABLE main.card_transactions
(
    id               SERIAL PRIMARY KEY,                 -- Идентификатор операции
    card_id          INTEGER REFERENCES main.cards (id), -- Внешний ключ на карту
    amount           DECIMAL(15, 2) NOT NULL,            -- Сумма операции
    transaction_type VARCHAR(50),                        -- Тип операции (оплата, снятие и т.д.)
    transaction_date TIMESTAMP DEFAULT NOW(),            -- Дата операции
    status           VARCHAR(50)                         -- Статус операции (успешно, отклонено и т.д.)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS main.card_transactions;
-- +goose StatementEnd
