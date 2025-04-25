-- +goose Up
CREATE TABLE products (
                          id SERIAL PRIMARY KEY,
                          name TEXT NOT NULL,
                          description TEXT,
                          price NUMERIC(10,2) NOT NULL,
                          availability BOOLEAN NOT NULL DEFAULT TRUE,
                          category_id INT REFERENCES categories(id) ON DELETE SET NULL,
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS products;
