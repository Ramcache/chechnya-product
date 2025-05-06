-- +goose Up
ALTER TABLE orders
    ADD COLUMN delivery_fee numeric(10, 2) DEFAULT 0,
    ADD COLUMN delivery_text text;

-- +goose Down
ALTER TABLE orders
    DROP COLUMN delivery_fee,
    DROP COLUMN delivery_text;
