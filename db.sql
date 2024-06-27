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

CREATE INDEX IF NOT EXISTS user_id_storage_asc ON orders (user_id, storage_until ASC);
CREATE INDEX IF NOT EXISTS id_asc ON orders (id ASC);