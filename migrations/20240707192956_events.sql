-- +goose Up
-- +goose StatementBegin
CREATE TYPE event_status AS ENUM ('none', 'acquired', 'processed');

CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    request TEXT NOT NULL,
    method_name TEXT NOT NULL,
    status event_status DEFAULT 'none',
    acquired_at TIMESTAMPTZ,
    processed_at TIMESTAMPTZ
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;
DROP TYPE event_status;
-- +goose StatementEnd