-- +goose Up
CREATE TABLE cart_items (
                            id SERIAL PRIMARY KEY,
                            cart_id INT REFERENCES carts(id),
                            product_id INT REFERENCES products(id),
                            quantity INT NOT NULL,
                            UNIQUE (cart_id, product_id)
);



-- +goose Down
DROP TABLE IF EXISTS cart_items;
