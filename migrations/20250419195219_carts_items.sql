-- +goose Up
CREATE TABLE cart_items (
                            id SERIAL PRIMARY KEY,
                            owner_id TEXT NOT NULL,                             -- например: user_1, guest_abc123
                            product_id INTEGER NOT NULL,
                            quantity INTEGER NOT NULL DEFAULT 1,
                            UNIQUE(owner_id, product_id),
                            FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS cart_items;
