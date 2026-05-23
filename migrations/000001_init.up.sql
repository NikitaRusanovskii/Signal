CREATE TABLE IF NOT EXISTS peers (
    id              UUID PRIMARY KEY,
    role            TEXT NOT NULL,
    is_online       BOOLEAN   DEFAULT false,
    addr_port       TEXT NOT NULL,
    connection_time TIMESTAMP NOT NULL
);