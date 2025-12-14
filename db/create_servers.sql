CREATE TABLE IF NOT EXISTS servers
(
    id           TEXT UNIQUE PRIMARY KEY,
    name         TEXT      NOT NULL,
    date_created TIMESTAMP NOT NULL
);