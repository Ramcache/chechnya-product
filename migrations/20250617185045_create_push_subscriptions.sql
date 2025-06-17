-- +goose Up
CREATE TABLE push_subscriptions (
                                    id SERIAL PRIMARY KEY,
                                    endpoint TEXT NOT NULL,
                                    p256dh TEXT NOT NULL,
                                    auth TEXT NOT NULL,
                                    user_id INT REFERENCES users(id) ON DELETE CASCADE,
                                    created_at TIMESTAMP DEFAULT now()
);

-- +goose Down
DROP TABLE push_subscriptions;
