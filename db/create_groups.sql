CREATE TABLE IF NOT EXISTS groups
(
    id           TEXT PRIMARY KEY REFERENCES channels (id),
    group_name   TEXT      NOT NULL UNIQUE,
    date_created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
