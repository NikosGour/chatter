CREATE TABLE IF NOT EXISTS servers
(
    id           TEXT UNIQUE PRIMARY KEY,
    "name"       TEXT      NOT NULL,
    date_created TIMESTAMP NOT NULL,
    is_test      boolean DEFAULT FALSE,

    CONSTRAINT name_it_test_unique UNIQUE ("name", is_test)
);
