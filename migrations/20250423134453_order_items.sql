-- +goose Up
CREATE TABLE IF NOT EXISTS order_items (
                                           id SERIAL PRIMARY KEY,
                                           order_id INT REFERENCES orders(id) ON DELETE CASCADE,
                                           product_id INT REFERENCES products(id) ON DELETE CASCADE,
                                           quantity INT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS order_items;
