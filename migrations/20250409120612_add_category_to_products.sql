-- +goose Up
ALTER TABLE products ADD COLUMN category TEXT DEFAULT '';

-- +goose Down
ALTER TABLE products DROP COLUMN category;
