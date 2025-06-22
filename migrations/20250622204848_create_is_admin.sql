-- +goose Up
ALTER TABLE push_subscriptions ADD COLUMN is_admin BOOLEAN DEFAULT false;

-- +goose Down
ALTER TABLE push_subscriptions DROP COLUMN is_admin;