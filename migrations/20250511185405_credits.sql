-- +goose Up
-- +goose StatementBegin
CREATE TABLE main.credits
(
    id             SERIAL PRIMARY KEY,                 -- Идентификатор кредита
    user_id        INTEGER REFERENCES main.users (id), -- Внешний ключ на пользователя
    amount         DECIMAL(15, 2) NOT NULL,            -- Сумма кредита
    interest_rate  DECIMAL(5, 2)  NOT NULL,            -- Процентная ставка
    term_in_months INTEGER        NOT NULL,            -- Срок кредита (в месяцах)
    status         VARCHAR(50),                        -- Статус кредита (оформлен, погашен и т.д.)
    created_at     TIMESTAMP DEFAULT NOW(),            -- Дата оформления кредита
    updated_at     TIMESTAMP DEFAULT NOW()             -- Дата последнего обновления
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS main.credits;
-- +goose StatementEnd
