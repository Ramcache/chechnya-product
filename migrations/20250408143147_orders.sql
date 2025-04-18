-- +goose Up
CREATE TABLE orders (
                        id SERIAL PRIMARY KEY,
                        owner_id TEXT NOT NULL,
                        total NUMERIC(10,2),
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS orders;
