-- +goose Up
-- +goose StatementBegin
CREATE TABLE main.transfers
(
    id              SERIAL PRIMARY KEY,                    -- Идентификатор перевода
    from_account_id INTEGER REFERENCES main.accounts (id), -- Внешний ключ на отправляющий счет
    to_account_id   INTEGER REFERENCES main.accounts (id), -- Внешний ключ на получающий счет
    amount          DECIMAL(15, 2) NOT NULL,               -- Сумма перевода
    transfer_date   TIMESTAMP DEFAULT NOW(),               -- Дата перевода
    status          VARCHAR(50)                            -- Статус перевода (успешно, отклонено)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS main.transfers;
-- +goose StatementEnd
