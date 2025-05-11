-- +goose Up
-- +goose StatementBegin
CREATE TABLE main.notifications
(
    id         SERIAL PRIMARY KEY,            -- Идентификатор уведомления
    user_id    INTEGER REFERENCES users (id), -- Внешний ключ на пользователя
    type       VARCHAR(50),                   -- Тип уведомления (email, SMS и т.д.)
    subject    VARCHAR(255),                  -- Тема письма
    body       TEXT,                          -- Текст письма
    status     VARCHAR(50),                   -- Статус отправки (отправлено, не отправлено)
    created_at TIMESTAMP DEFAULT NOW()        -- Дата отправки
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS main.notifications;
-- +goose StatementEnd
