-- +goose Up
ALTER TABLE orders
    ADD COLUMN comment TEXT,
    ADD COLUMN rating INTEGER CHECK (rating BETWEEN 1 AND 5);

-- +goose Down
ALTER TABLE orders
    DROP COLUMN rating,
    DROP COLUMN comment;
