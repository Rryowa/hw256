CREATE TABLE orders (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    storage_until TIMESTAMPTZ NOT NULL,
    issued BOOLEAN NOT NULL,
    issued_at TIMESTAMPTZ,
    returned BOOLEAN NOT NULL,
    hash VARCHAR(255) NOT NULL
);

CREATE INDEX idx_orders_returned ON orders(returned);
CREATE INDEX idx_orders_user_id_issued ON orders(user_id, issued);