-- +goose Up
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username VARCHAR(100),
                       email VARCHAR(255) UNIQUE,
                       phone VARCHAR(20) NOT NULL UNIQUE,
                       password_hash VARCHAR(255) NOT NULL,
                       role VARCHAR(50) DEFAULT 'user',
                       is_verified BOOLEAN DEFAULT TRUE,
                       owner_id VARCHAR(100) UNIQUE,
                       created_at TIMESTAMP DEFAULT NOW()
);



-- +goose Down
DROP TABLE IF EXISTS users;
