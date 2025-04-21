-- +goose Up
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username VARCHAR(100),                         -- можно оставить уникальность по желанию
                       email VARCHAR(255) UNIQUE,                     -- nullable
                       phone VARCHAR(20) NOT NULL UNIQUE,             -- регистрация по телефону
                       password_hash VARCHAR(255) NOT NULL,
                       role VARCHAR(50) DEFAULT 'user',
                       is_verified BOOLEAN DEFAULT FALSE,
                       owner_id VARCHAR(100) UNIQUE,                  -- user_xxx или guest_xxx
                       created_at TIMESTAMP DEFAULT NOW()
);



-- +goose Down
DROP TABLE IF EXISTS users;
