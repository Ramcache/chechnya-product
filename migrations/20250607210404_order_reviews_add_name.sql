-- +goose Up
ALTER TABLE order_reviews
    ADD COLUMN user_id INT REFERENCES users(id) ON DELETE SET NULL;

-- +goose Down
ALTER TABLE order_reviews
    DROP COLUMN user_id;
