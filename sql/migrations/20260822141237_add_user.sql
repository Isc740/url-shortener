-- +goose Up
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    user_name TEXT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    exists BOOL DEFAULT TRUE
);

-- +goose Down
DROP TABLE users;
