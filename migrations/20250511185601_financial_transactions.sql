-- +goose Up
-- +goose StatementBegin
CREATE TABLE main.financial_transactions
(
    id                 SERIAL PRIMARY KEY,                 -- Идентификатор операции
    user_id            INTEGER REFERENCES main.users (id), -- Внешний ключ на пользователя
    transaction_type   VARCHAR(50),                        -- Тип операции (пополнение, снятие, перевод и т.д.)
    amount             DECIMAL(15, 2) NOT NULL,            -- Сумма операции
    transaction_date   TIMESTAMP DEFAULT NOW(),            -- Дата операции
    transaction_status VARCHAR(50),                        -- Статус операции (успешно, отклонено)
    created_at         TIMESTAMP DEFAULT NOW(),            -- Дата создания записи
    updated_at         TIMESTAMP DEFAULT NOW(),            -- Дата последнего обновления
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS main.financial_transactions;
-- +goose StatementEnd
