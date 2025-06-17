-- +goose Up
ALTER TABLE orders
    ADD COLUMN latitude DOUBLE PRECISION,
    ADD COLUMN longitude DOUBLE PRECISION;

-- +goose Down
ALTER TABLE orders
    DROP COLUMN latitude,
    DROP COLUMN longitude;