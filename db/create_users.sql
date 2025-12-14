CREATE TABLE IF NOT EXISTS users
(
    id           TEXT UNIQUE PRIMARY KEY,
    username     TEXT UNIQUE NOT NULL,
    password     TEXT        NOT NULL,
    date_created TIMESTAMP   NOT NULL
);