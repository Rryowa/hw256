-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    storage_until TIMESTAMPTZ NOT NULL,
    issued BOOLEAN NOT NULL,
    issued_at TIMESTAMPTZ,
    returned BOOLEAN NOT NULL,
    order_price FLOAT NOT NULL,
    weight FLOAT NOT NULL,
    package_type VARCHAR(255) NOT NULL,
    package_price FLOAT NOT NULL,
    hash VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    request TEXT NOT NULL,
    method_name TEXT NOT NULL,
    acquired BOOLEAN DEFAULT FALSE,
    processed BOOLEAN DEFAULT FALSE,
    acquired_at TIMESTAMPTZ,
    processed_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS user_id_storage_asc ON orders (user_id, storage_until ASC);
CREATE INDEX IF NOT EXISTS id_asc ON orders (id ASC);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS user_id_storage_asc;
DROP INDEX IF EXISTS id_asc;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS events;
-- +goose StatementEnd