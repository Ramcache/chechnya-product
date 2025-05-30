-- +goose Up
CREATE TABLE order_reviews (
                               id SERIAL PRIMARY KEY,
                               order_id INT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
                               rating INTEGER CHECK (rating BETWEEN 1 AND 5),
                               comment TEXT,
                               created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE order_reviews;
