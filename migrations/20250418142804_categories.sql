-- +goose Up
CREATE TABLE categories (
                            id SERIAL PRIMARY KEY,
                            name TEXT UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS categories;
