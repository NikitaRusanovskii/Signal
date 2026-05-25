CREATE TABLE IF NOT EXISTS peers (
    role            TEXT NOT NULL,
    is_online       BOOLEAN   DEFAULT true,
    addr_port       TEXT NOT NULL UNIQUE,
    last_seen TIMESTAMP NOT NULL
);
