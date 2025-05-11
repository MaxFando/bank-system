-- +goose Up
-- +goose StatementBegin
CREATE TABLE main.payment_schedules
(
    id               SERIAL PRIMARY KEY,                   -- Идентификатор записи
    credit_id        INTEGER REFERENCES main.credits (id), -- Внешний ключ на кредит
    payment_date     TIMESTAMP,                            -- Дата платежа
    payment_amount   DECIMAL(15, 2) NOT NULL,              -- Сумма платежа
    principal_amount DECIMAL(15, 2) NOT NULL,              -- Основной долг
    interest_amount  DECIMAL(15, 2) NOT NULL,              -- Сумма процентов
    penalty          DECIMAL(10, 2) DEFAULT 0,             -- Штраф за просрочку
    balance          DECIMAL(10, 2),                       -- Остаток долга после платежа
    created_at       TIMESTAMP      DEFAULT NOW(),         -- Дата создания записи
    updated_at       TIMESTAMP      DEFAULT NOW()          -- Дата последнего обновления
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS main.payment_schedules;
-- +goose StatementEnd
