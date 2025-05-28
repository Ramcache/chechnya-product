-- +goose Up
ALTER TABLE orders
    ADD COLUMN order_comment TEXT;

-- +goose Down
ALTER TABLE orders
    DROP COLUMN order_comment;
