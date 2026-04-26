-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id         UUID        PRIMARY KEY,
    name       TEXT        NOT NULL,
    email      TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX users_email_idx
    ON users (email)
    WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS users_email_idx;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
