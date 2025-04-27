-- +goose Up
ALTER TABLE products
    ADD COLUMN url TEXT;

-- +goose Down
ALTER TABLE products
    DROP COLUMN url;
