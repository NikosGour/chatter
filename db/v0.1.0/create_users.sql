CREATE TABLE IF NOT EXISTS users
(
    id           TEXT PRIMARY KEY REFERENCES channels (id),
    username     TEXT      NOT NULL UNIQUE,
    "password"   TEXT      NOT NULL,
    date_created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
