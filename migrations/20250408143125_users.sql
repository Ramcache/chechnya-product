-- +goose Up
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username VARCHAR(100) NOT NULL UNIQUE,      -- имя обязательно и уникально
                       email VARCHAR(255) UNIQUE,                  -- email уникальный, может быть NULL
                       phone VARCHAR(20) NOT NULL UNIQUE,          -- телефон обязательно и уникально
                       password_hash VARCHAR(255) NOT NULL,        -- пароль обязательно
                       role VARCHAR(50) NOT NULL DEFAULT 'user',   -- роль обязательно, по умолчанию user
                       is_verified BOOLEAN NOT NULL DEFAULT TRUE,  -- обязательно, по умолчанию TRUE
                       owner_id VARCHAR(100) NOT NULL UNIQUE,      -- обязательно и уникально
                       created_at TIMESTAMP NOT NULL DEFAULT NOW() -- обязательно, по умолчанию текущее время
);

-- Индекс на owner_id для быстрого поиска (уже уникальный, но индекс ускоряет выборку)
CREATE INDEX idx_users_owner_id ON users(owner_id);

-- Индекс на username для поиска (уже уникальный, но индекс ускоряет выборку)
CREATE INDEX idx_users_username ON users(username);

-- +goose Down
DROP INDEX IF EXISTS idx_users_owner_id;
DROP INDEX IF EXISTS idx_users_username;
DROP TABLE IF EXISTS users;
