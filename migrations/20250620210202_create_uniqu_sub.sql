-- +goose Up
CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_subscription ON push_subscriptions(endpoint);

-- +goose Down
DROP INDEX IF EXISTS idx_unique_subscription;
