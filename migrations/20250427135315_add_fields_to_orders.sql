-- +goose Up
ALTER TABLE orders
    ADD COLUMN name TEXT,
    ADD COLUMN address TEXT,
    ADD COLUMN delivery_type TEXT,
    ADD COLUMN payment_type TEXT,
    ADD COLUMN change_for NUMERIC(10,2);

-- +goose Down
ALTER TABLE orders
    DROP COLUMN name,
    DROP COLUMN address,
    DROP COLUMN delivery_type,
    DROP COLUMN payment_type,
    DROP COLUMN change_for;
