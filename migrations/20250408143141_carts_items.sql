-- +goose Up
CREATE TABLE cart_items (
                            id SERIAL PRIMARY KEY,
                            cart_id INT REFERENCES carts(id) ON DELETE CASCADE,
                            product_id INT REFERENCES products(id),
                            quantity INT NOT NULL
);
-- +goose Down
DROP TABLE IF EXISTS cart_items;
