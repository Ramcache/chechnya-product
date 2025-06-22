-- +goose Up
CREATE TABLE IF NOT EXISTS push_subscriptions (
                                                  id SERIAL PRIMARY KEY,
                                                  endpoint TEXT NOT NULL UNIQUE,
                                                  p256dh TEXT NOT NULL,
                                                  auth TEXT NOT NULL,
                                                  created_at TIMESTAMP DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS push_subscriptions;
