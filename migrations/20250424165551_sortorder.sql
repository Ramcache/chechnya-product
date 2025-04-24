-- +goose Up
ALTER TABLE categories ADD COLUMN sort_order INT DEFAULT 0;

-- -goose Down
ALTER TABLE categories DROP COLUMN sort_order;
