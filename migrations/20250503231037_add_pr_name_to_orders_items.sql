-- +goose Up
ALTER TABLE order_items
    ADD COLUMN product_name varchar(255),
    ADD COLUMN price numeric(10, 2);

-- +goose Down
ALTER TABLE order_items
    DROP COLUMN product_name,
    DROP COLUMN price