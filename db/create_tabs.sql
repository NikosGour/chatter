CREATE TABLE IF NOT EXISTS tabs
(
    id           TEXT UNIQUE PRIMARY KEY,
    "name"       TEXT      NOT NULL,
    server_id    TEXT      NOT NULL REFERENCES servers (id) ON DELETE CASCADE,
    date_created TIMESTAMP NOT NULL
);