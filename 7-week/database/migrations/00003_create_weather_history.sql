-- +goose Up
-- +goose StatementBegin
CREATE TABLE weather_history (
    id           UUID         PRIMARY KEY,
    user_id      UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    city         TEXT         NOT NULL,
    temperature  NUMERIC(5,2) NOT NULL,
    description  TEXT         NOT NULL,
    requested_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX weather_history_user_id_city_idx
    ON weather_history (user_id, city);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS weather_history_user_id_city_idx;
DROP TABLE IF EXISTS weather_history;
-- +goose StatementEnd
