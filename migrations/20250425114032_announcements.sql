-- +goose Up
CREATE TABLE announcements (
                               id SERIAL PRIMARY KEY,
                               title TEXT NOT NULL,
                               content TEXT NOT NULL,
                               created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS announcements;
