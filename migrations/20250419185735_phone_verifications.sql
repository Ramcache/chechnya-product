-- +goose Up
CREATE TABLE phone_verifications (
                                     phone VARCHAR PRIMARY KEY,
                                     code VARCHAR NOT NULL,
                                     created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                     confirmed BOOLEAN NOT NULL DEFAULT FALSE
);

-- +goose Down
DROP TABLE IF EXISTS phone_verifications;
