-- +goose Up
CREATE TABLE verification_codes (
                                    phone TEXT PRIMARY KEY,
                                    code TEXT NOT NULL,
                                    expires_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE verification_codes;