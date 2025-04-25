-- +goose Up
CREATE TABLE categories (
                            id SERIAL PRIMARY KEY,
                            name TEXT UNIQUE NOT NULL,
                            sort_order INT DEFAULT 0
);

-- +goose Down
DROP TABLE IF EXISTS categories;
