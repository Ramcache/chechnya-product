-- +goose Up
CREATE TABLE reviews (
                         id SERIAL PRIMARY KEY,
                         owner_id TEXT NOT NULL,
                         product_id INT NOT NULL,
                         rating INT CHECK (rating >= 1 AND rating <= 5),
                         comment TEXT,
                         created_at TIMESTAMP DEFAULT NOW(),

                         CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);


-- +goose Down
DROP TABLE IF EXISTS reviews;