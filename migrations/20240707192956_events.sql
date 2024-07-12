-- +goose Up
-- +goose StatementBegin
CREATE TYPE event_status AS ENUM ('none', 'requested');

CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    method_name TEXT NOT NULL,
    request TEXT NOT NULL,
    status event_status DEFAULT 'none',
    requested_at TIMESTAMPTZ
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;
DROP TYPE event_status;
-- +goose StatementEnd