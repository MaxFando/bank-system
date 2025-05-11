-- +goose Up
-- +goose StatementBegin
CREATE TABLE main.central_bank_rates
(
    id         SERIAL PRIMARY KEY,     -- Идентификатор записи
    rate       DECIMAL(5, 2) NOT NULL, -- Ключевая ставка
    rate_date  DATE          NOT NULL, -- Дата ставки
    created_at TIMESTAMP DEFAULT NOW() -- Дата записи
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS main.central_bank_rates;
-- +goose StatementEnd
