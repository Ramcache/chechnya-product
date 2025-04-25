-- +goose Up
ALTER TABLE orders ADD COLUMN status TEXT NOT NULL DEFAULT 'в обработке';

-- +goose Down
ALTER TABLE orders DROP COLUMN status;