-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    storage_until TIMESTAMPTZ NOT NULL,
    issued BOOLEAN NOT NULL,
    issued_at TIMESTAMPTZ,
    returned BOOLEAN NOT NULL,
    hash VARCHAR(255) NOT NULL
);


-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE orders;
-- +goose StatementEnd