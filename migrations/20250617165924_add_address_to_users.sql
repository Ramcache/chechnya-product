-- +goose Up
ALTER TABLE users ADD COLUMN address TEXT;

-- +goose Down
ALTER TABLE users DROP COLUMN address;
