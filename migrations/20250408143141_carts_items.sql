-- +goose Up
CREATE TABLE cart_items (
                            id SERIAL PRIMARY KEY,
                            owner_id TEXT NOT NULL,                        -- универсальный идентификатор: user_1 или ip_192.168.1.1
                            product_id INTEGER NOT NULL,
                            quantity INTEGER NOT NULL DEFAULT 1,
                            UNIQUE(owner_id, product_id),
                            FOREIGN KEY (product_id) REFERENCES products(id)
);

-- +goose Down
DROP TABLE IF EXISTS cart_items;
