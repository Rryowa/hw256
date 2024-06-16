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

-- CREATE INDEX user_id_hash ON orders using hash(user_id);
-- CREATE INDEX storage_until_b_tree ON orders (storage_until DESC);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
-- DROP INDEX user_id_hash;
-- DROP INDEX storage_until_b_tree;
DROP TABLE orders;
-- +goose StatementEnd